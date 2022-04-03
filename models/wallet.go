package models

import (
	"errors"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ID     int32  `db:"id, primarykey, autoincrement" json:"id"`
	UserID int32  `db:"user_id" json:"-"`
	Title  string `db:"title" json:"title"`
}

func IsWalletAvailable(db *gorm.DB, walletID int32) bool {
	var wallet Wallet
	err := db.First(&wallet, walletID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
