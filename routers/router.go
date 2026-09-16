package routers

import (
	"github.com/arunkumar-1311/mongo-db/handlers"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.New()

	v1 := r.Group("v1/api")
	v1.POST("/user/create", handlers.CreateUsers)

	r.Run(":8000")
}
