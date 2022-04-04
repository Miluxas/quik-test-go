package controllers

import (
	"net/http"

	"miluxas/models"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

var authModel = new(models.AuthModel)

func (ctl AuthController) TokenValid(context *gin.Context) {

	UserID, err := authModel.ExtractTokenMetadata(context.Request)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}
	context.Set("userID", UserID)
}
