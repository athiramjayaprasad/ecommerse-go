package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/athiramjayaprasad/ecommerse-go/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type SignedDetails struct {
	Email string
	First_Name string
	Last_Name string
	Uid string
	jwt.StandardClaims

}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator( email string, first_name string, last_name string, uid string)(signed_token string, signed_refresh_token string, err error)  {
	claims := &SignedDetails{
		Email: email,
		First_Name: first_name,
		Last_Name: last_name,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour *time.Duration(24)).Unix(),
		},
	}
	refresh_claims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodES384, refresh_claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refresh_token, err
}

func ValidateToken(signed_token string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signed_token, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if  err !=nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return 
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = "token is already expired"
		return
	}
	return claims, msg
}
	 

func UpdateAllTokens(signed_token string, signed_refresh_token string, user_id string)  {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var update_obj primitive.D
	update_obj = append(update_obj, bson.E{Key: "token", Value: signed_token})
	update_obj = append(update_obj, bson.E{Key: "refresh_token", Value: signed_refresh_token})
	upadated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update_obj = append(update_obj, bson.E{Key: "updated_at", Value: upadated_at})
	upsert := true
	filter := bson.M{"user_id":user_id}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: update_obj},
	}, &opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}
