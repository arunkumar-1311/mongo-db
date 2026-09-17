package routers

import (
	"github.com/arunkumar-1311/mongo-db/handlers"
	"github.com/arunkumar-1311/mongo-db/middleware"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.New()

	v1 := r.Group("v1/api")

	// Public routes.
	v1.POST("/user/create", handlers.CreateUsers)
	v1.POST("/auth/login", handlers.Login)

	// Authenticated routes.
	authorized := v1.Group("")
	authorized.Use(middleware.Authenticate())
	authorized.GET("/user/get/:id", handlers.GetUser)

	r.Run(":8000")
}
