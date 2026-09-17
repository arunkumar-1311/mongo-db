package handlers

import (
	"net/http"

	"github.com/arunkumar-1311/mongo-db/models"
	"github.com/arunkumar-1311/mongo-db/service"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "unable to decode the request body: " + err.Error()})
		return
	}

	result, err := service.NewUserService().Login(c.Request.Context(), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid id: " + err.Error()})
		return
	}

	result, err := service.NewUserService().GetUser(c.Request.Context(), objID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "unable to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
