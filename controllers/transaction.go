package controllers

import (
	"strconv"

	"miluxas/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type TransactionController struct{}

var transactionModel = new(models.TransactionModel)

type CreditRequestBody struct {
	Value decimal.Decimal
}

func (ctrl TransactionController) Credit(context *gin.Context) {
	walletId, _ := context.Params.Get("wallet_id")
	var requestBody CreditRequestBody
	context.BindJSON(&requestBody)
	if !requestBody.Value.IsPositive() {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Value should be a positive decimal"})
		return
	}
	walletIdInt, _ := strconv.ParseInt(walletId, 10, 32)
	newBalance, err := transactionModel.RegisterCredit(Db, int32(walletIdInt), requestBody.Value)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Credit transaction registered successfully", "newBalance": newBalance})
}

type DebitRequestBody struct {
	Value decimal.Decimal
}

func (ctrl TransactionController) Debit(context *gin.Context) {
	walletId, _ := context.Params.Get("wallet_id")
	var requestBody DebitRequestBody
	context.BindJSON(&requestBody)
	if !requestBody.Value.IsPositive() {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Value should be a positive decimal"})
		return
	}
	walletIdInt, _ := strconv.ParseInt(walletId, 10, 32)
	newBalance, err := transactionModel.RegisterDebit(Db, int32(walletIdInt), requestBody.Value)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Debit transaction registered successfully", "newBalance": newBalance})
}

func (ctrl TransactionController) Balance(context *gin.Context) {
	walletId, _ := context.Params.Get("wallet_id")
	walletIdInt, _ := strconv.ParseInt(walletId, 10, 32)
	newBalance, err := models.GetBalance(Db, int32(walletIdInt))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "", "Balance": newBalance})
}
