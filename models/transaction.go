package models

import (
	"errors"
	Db "miluxas/db"

	log "github.com/sirupsen/logrus"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID       int32           `db:"id, primarykey, autoincrement" json:"id"`
	WalletID int32           `db:"wallet_id" json:"-"`
	Credit   decimal.Decimal `db:"credit" json:"credit"`
	Debit    decimal.Decimal `db:"debit" json:"debit"`
}

type TransactionModel struct{}

func (m TransactionModel) RegisterCredit(db *gorm.DB, walletID int32, value decimal.Decimal) (newBalance decimal.Decimal, err error) {

	isValidWallet := IsWalletAvailable(db, walletID)
	if !isValidWallet {
		return decimal.Zero.Floor(), errors.New("Wallet is not available")
	}

	var newTransaction Transaction
	newTransaction.WalletID = walletID
	newTransaction.Credit = value
	db.Create(&newTransaction)
	newBalance, _ = getBalance(db, walletID)
	Db.SetCacheBalance(walletID, newBalance)

	return newBalance, nil
}

func (m TransactionModel) RegisterDebit(db *gorm.DB, walletID int32, value decimal.Decimal) (newBalance decimal.Decimal, err error) {

	isValidWallet := IsWalletAvailable(db, walletID)
	if !isValidWallet {
		return decimal.Zero.Floor(), errors.New("Wallet is not available")
	}

	var newTransaction Transaction
	newTransaction.WalletID = walletID
	newTransaction.Debit = value
	newBalance, _ = GetBalance(db, walletID)
	if newBalance.Cmp(value) == -1 {
		return newBalance, errors.New("Debit value id mode than balance")
	}
	db.Create(&newTransaction)
	newBalance, _ = getBalance(db, walletID)
	Db.SetCacheBalance(walletID, newBalance)

	return newBalance, nil
}

func getBalance(db *gorm.DB, walletID int32) (balance decimal.Decimal, err error) {
	isValidWallet := IsWalletAvailable(db, walletID)
	if !isValidWallet {
		return decimal.Zero.Floor(), errors.New("Wallet is not available")
	}
	err = db.Table("transactions").Where("wallet_id = ?", walletID).Select("sum(credit-debit)").Row().Scan(&balance)
	if err != nil {
		return decimal.Zero.Floor(), errors.New("Something went wrong")
	}
	return balance, nil
}

func GetBalance(db *gorm.DB, walletID int32) (balance decimal.Decimal, err error) {

	isValidWallet := IsWalletAvailable(db, walletID)
	if !isValidWallet {
		return decimal.Zero.Floor(), errors.New("Wallet is not available")
	}
	cacheBalance, err := Db.GetCacheBalance(walletID)
	if err != nil {
		newBalance, _ := getBalance(db, walletID)
		Db.SetCacheBalance(walletID, newBalance)
		return newBalance, nil
	} else {
		log.Infoln("Read wallet balance from cache")
		return cacheBalance, nil
	}
}
