package router

import (
	"github.com/Goryudyuma/kstmisucon1/controllers"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	//	r.GET("/", controllers.APIEndpoints)

	api := r.Group("api")
	{

		api.GET("/comments", controllers.GetComments)
		api.GET("/comments/:id", controllers.GetComment)
		api.POST("/comments", controllers.CreateComment)
		api.PUT("/comments/:id", controllers.UpdateComment)
		api.DELETE("/comments/:id", controllers.DeleteComment)

		api.GET("/users", controllers.GetUsers)
		api.GET("/users/:id", controllers.GetUser)
		api.POST("/users", controllers.CreateUser)
		api.PUT("/users/:id", controllers.UpdateUser)
		api.DELETE("/users/:id", controllers.DeleteUser)

	}
}
