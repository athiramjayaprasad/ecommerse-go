package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/athiramjayaprasad/ecommerse-go/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine)  {
	incomingRoutes.POST("/users/sign_up", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/add_product", controllers.ProductViewerAdmin())
	incomingRoutes.POST("/users/product_view", controllers.SearchProduct())
	incomingRoutes.POST("/users/search", controllers.SearchProductByQuery())  
}                                                                        