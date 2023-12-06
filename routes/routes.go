// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"

	"gin-api/controllers"
)

// InitRoutes initializes the routes
func InitRoutes(router *gin.Engine) {
	v1 := router.Group("/api/users")
	{
		v1.POST("/create", controllers.CreateUser)
		v1.GET("/", controllers.GetUsers)
		v1.GET("/:id", controllers.GetUserByID)
		v1.PUT("/update/:id", controllers.UpdateUser)
		v1.DELETE("/delete/:id", controllers.DeleteUser)
	}
}
