package controllers

import (
	"miluxas/models"

	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

var Db *gorm.DB
var userModel = new(models.UserModel)

func getUserID(c *gin.Context) (userID int64) {
	return c.MustGet("userID").(int64)
}

type LoginRequestBody struct {
	Email    string
	Password string
}

func (ctrl UserController) Login(context *gin.Context) {
	var RequestBody LoginRequestBody
	context.BindJSON(&RequestBody)

	user, token, err := userModel.Login(Db, RequestBody.Email, RequestBody.Password)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid login details"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}
