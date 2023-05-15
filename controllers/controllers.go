package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/athiramjayaprasad/ecommerse-go/database"
	"github.com/athiramjayaprasad/ecommerse-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()    
    
        
func HashPassword(password string) string  {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}  
	return string(bytes)                          
}
  
func VerifyPassword (userPassword string, givenPassword string)(bool, string) {
	 err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	 valid := true
	 msg := ""
	 
	 if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	 }
	 return valid, msg 
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx , cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}  

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} 

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"user already exists"}) 
		}  
		
		UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"this phone number is already in use"}) 
		}  

		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()

		token, refreshtoken, _ := generate.TokenGenerator                                                                   
		user.Token 
		user.Refresh_Token
		user.User_Cart
		user.Address_Details
		user.Order_Status
	}  
}
func ProductViewerAdmin() gin.HandlerFunc { 
	                                    
}
func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong, please try after some time")
			return
		}
		err = cursor.All(ctx, &productList)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return  
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")           
		}
		c.IndentedJSON(200, productList)
		
	}                                        
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}}) 

		if err != nil {
			c.IndentedJSON(404, "Something went wrong while fetching the data")
			return
		}

		err = searchquerydb.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}  

		defer searchquerydb.Close(ctx)
		
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchProducts)                                                     
	}                                         
	
}
