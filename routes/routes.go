// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"

	"gin-api/controllers"
	"gin-api/helpers"
	"gin-api/middleware"
)

// InitRoutes initializes the routes
func InitRoutes(router *gin.Engine) {
	users := router.Group("/api/users")
	{
		users.POST("/create", controllers.CreateUser)
		users.GET("/", controllers.GetUsers)
		users.GET("/:id", controllers.GetUserByID)
		users.PUT("/update/:id", controllers.UpdateUser)
		users.DELETE("/delete/:id", controllers.DeleteUser)
		users.GET("/test/:id", controllers.OneUsersHandler)
	}

	auth := router.Group("/api/auth")
	{
		auth.POST("/signin", controllers.Login)
	}

	roles := router.Group("/api/roles")
	{
		roles.POST("/create", controllers.CreateRole)
		roles.GET("/allRoles", controllers.GetAllRoles)
	}

	product := router.Group("/api/product")
	{
		product.POST("/createProduct", middleware.EnsureAdmin(), controllers.CreateProduct)
		product.POST("/createTransProduct", middleware.EnsureAdmin(), controllers.CreateProduct)
		product.GET("/allProduct", middleware.EnsureAdmin(), controllers.AllProduct)
		product.GET("/oneProduct/:id", middleware.EnsureAdmin(), controllers.OneProduct)
		product.PUT("/updateProduct/:id", middleware.EnsureAdmin(), controllers.UpdateProduct)
		product.DELETE("/deleteProduct/:id", middleware.EnsureAdmin(), controllers.DeleteProduct)
		product.GET("/image/:filename", helpers.ShowImageFromMinio)
		product.GET("/download/:filename", helpers.DownloadImage)
	}

	categori := router.Group("/api/categori")
	{
		categori.POST("/createCategori", middleware.EnsureAdmin(), controllers.CreateCategori)
		categori.GET("/allCategori", middleware.EnsureAdmin(), controllers.AllCategories)
		categori.GET("/oneCategori/:id", middleware.EnsureAdmin(), controllers.OneCategori)
	}
}
