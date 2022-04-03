package main

import (
	"fmt"
	"log"
	"miluxas/controllers"
	"miluxas/db"
	"os"

	"github.com/gin-contrib/gzip"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Content-Type", "application/json")

		if context.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}

var auth = new(controllers.AuthController)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth.TokenValid(context)
		context.Next()
	}
}

func main() {
	loadEnv()

	serverEngine := gin.Default()

	serverEngine.Use(CORSMiddleware())
	serverEngine.Use(gzip.Gzip(gzip.DefaultCompression))

	db.Init()
	db.InitRedis()

	controllers.Db = db.GetDB()

	user := new(controllers.UserController)
	serverEngine.POST("/api/v1/user/login", user.Login)

	transaction := new(controllers.TransactionController)
	serverEngine.POST("/api/v1/wallets/:wallet_id/credit", TokenAuthMiddleware(), transaction.Credit)
	serverEngine.POST("/api/v1/wallets/:wallet_id/debit", TokenAuthMiddleware(), transaction.Debit)
	serverEngine.GET("/api/v1/wallets/:wallet_id/balance", TokenAuthMiddleware(), transaction.Balance)

	port := os.Getenv("PORT")

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	serverEngine.Run(":" + port)

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}
}
