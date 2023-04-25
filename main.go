package main

import (
	"log"
	"os"

	"github.com/athiramjayaprasad/ecommerse-go/controllers"
	"github.com/athiramjayaprasad/ecommerse-go/database"
	"github.com/athiramjayaprasad/ecommerse-go/middleware"
	"github.com/athiramjayaprasad/ecommerse-go/routes"
	"github.com/gin-gonic/gin"
)

func main()  {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"

	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/add_to_cart", app.AddToCart())
	router.GET("/remove_item", app.RemoveItem())
	router.GET("cart_checkout", app.BuyFromCart())
	router.GET("/instant_buy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))  
	           
}