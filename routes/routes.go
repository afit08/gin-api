// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"

	"gin-api/controllers"
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

	private := users.Group("/apit")
	private.Use(middleware.EnsureAdmin())
	{
		private.GET("/ok/:id", controllers.OneUsersHandler)
	}

	product := router.Group("/api/product")
	{
		product.POST("/createProduct", middleware.EnsureAdmin(), controllers.CreateProduct)
		product.GET("/allProduct", middleware.EnsureAdmin(), controllers.AllProduct)
	}
}
