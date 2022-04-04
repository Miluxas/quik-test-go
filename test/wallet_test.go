package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"miluxas/controllers"
	"miluxas/db"
	"os"

	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var auth = new(controllers.AuthController)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.TokenValid(c)
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	serverEngine := gin.Default()
	gin.SetMode(gin.TestMode)

	user := new(controllers.UserController)
	serverEngine.POST("/api/v1/user/login", user.Login)

	transaction := new(controllers.TransactionController)
	serverEngine.POST("/api/v1/wallets/:wallet_id/credit", TokenAuthMiddleware(), transaction.Credit)
	serverEngine.POST("/api/v1/wallets/:wallet_id/debit", TokenAuthMiddleware(), transaction.Debit)
	serverEngine.GET("/api/v1/wallets/:wallet_id/balance", TokenAuthMiddleware(), transaction.Balance)

	return serverEngine
}

func main() {
	r := SetupRouter()
	r.Run()
}

var loginCookie string

var testEmail = "test@t.co"
var testPassword = "123123"

var walletID = "2"

var accessToken string

var articleID int

func TestIntDB(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file, please create one in the root directory")
	}

	fmt.Println("DB_PASS", os.Getenv("DB_PASS"))

	db.Init()
	db.InitRedis()
	controllers.Db = db.GetDB()

}

func TestLoginInvalidEmail(t *testing.T) {
	testRouter := SetupRouter()

	var loginRequest controllers.LoginRequestBody

	loginRequest.Email = "invalid@email"
	loginRequest.Password = testPassword

	data, _ := json.Marshal(loginRequest)

	req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotAcceptable, resp.Code)
}

func TestLoginInvalidPassword(t *testing.T) {
	testRouter := SetupRouter()

	var loginRequest controllers.LoginRequestBody

	loginRequest.Email = testEmail
	loginRequest.Password = "wrongpassword"

	data, _ := json.Marshal(loginRequest)

	req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotAcceptable, resp.Code)
}
func TestLogin(t *testing.T) {
	testRouter := SetupRouter()

	var loginRequest controllers.LoginRequestBody

	loginRequest.Email = testEmail
	loginRequest.Password = testPassword

	data, _ := json.Marshal(loginRequest)

	req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		Message string `json:"message"`
		User    struct {
			CreatedAt int64  `json:"created_at"`
			Email     string `json:"email"`
			ID        int64  `json:"id"`
			Name      string `json:"name"`
			UpdatedAt int64  `json:"updated_at"`
		} `json:"user"`
		Token string `json:"token"`
	}
	json.Unmarshal(body, &res)

	accessToken = res.Token

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetBalance(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", "/api/v1/wallets/"+walletID+"/balance", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		message string
		balance string
	}
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetBalanceInvalidWalletId(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("GET", "/api/v1/wallets/100/balance", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		message string
		balance string
	}
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreditNegativeValue(t *testing.T) {
	testRouter := SetupRouter()

	var requestBody controllers.CreditRequestBody

	requestBody.Value, _ = decimal.NewFromString("-10")
	data, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/wallets/"+walletID+"/credit", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		message    string
		newBalance string
	}
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCredit(t *testing.T) {
	testRouter := SetupRouter()

	var requestBody controllers.CreditRequestBody

	requestBody.Value, _ = decimal.NewFromString("10")
	data, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/wallets/"+walletID+"/credit", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		message    string
		newBalance string
	}
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDebitNegativeValue(t *testing.T) {
	testRouter := SetupRouter()

	var requestBody controllers.DebitRequestBody

	requestBody.Value, _ = decimal.NewFromString("-10")
	data, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/wallets/"+walletID+"/debit", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		message    string
		newBalance string
	}
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDebitOverThanBalance(t *testing.T) {
	testRouter := SetupRouter()

	var requestBody controllers.DebitRequestBody

	requestBody.Value, _ = decimal.NewFromString("1000")
	data, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/wallets/"+walletID+"/debit", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		message    string
		newBalance string
	}
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDebit(t *testing.T) {
	testRouter := SetupRouter()

	var requestBody controllers.DebitRequestBody

	requestBody.Value, _ = decimal.NewFromString("3")
	data, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/v1/wallets/"+walletID+"/debit", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res struct {
		message    string
		newBalance string
	}
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusOK, resp.Code)
}
