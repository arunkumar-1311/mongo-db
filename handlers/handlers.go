package handlers

import (
	"net/http"

	"github.com/arunkumar-1311/mongo-db/models"
	"github.com/arunkumar-1311/mongo-db/service"

	"github.com/gin-gonic/gin"
)

func CreateUsers(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "unable to decode the request body: " + err.Error()})
		return
	}

	result, err := service.NewUserService().CreateUsers(c.Request.Context(), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "unable to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
