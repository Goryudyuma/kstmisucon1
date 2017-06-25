package router

import (
	"github.com/Goryudyuma/kstmisucon1/controllers"
	"github.com/Goryudyuma/kstmisucon1/sessions"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	//	r.GET("/", controllers.APIEndpoints)

	api := r.Group("api")
	{
		api.GET("/comments/:id", controllers.GetComment)
		api.POST("/comments", controllers.CreateComment)
		api.PUT("/comments/:id", controllers.UpdateComment)

		api.GET("/users/:id", controllers.GetUserById)
		api.POST("/users", controllers.CreateUser)
		api.PUT("/users/:id", controllers.UpdateUser)

		api.POST("/login", sessions.Login)
		api.POST("/logout", sessions.Logout)
	}

	debug := r.Group("debug")
	{
		debug.GET("/users", controllers.GetUsers)
		debug.GET("/comments", controllers.GetComments)
	}
}
