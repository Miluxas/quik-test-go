package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int32  `db:"id, primarykey, autoincrement" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
	Name     string `db:"name" json:"name"`
}

type UserModel struct{}

var authModel = new(AuthModel)

func (m UserModel) Login(db *gorm.DB, email string, password string) (user User, token string, err error) {
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, token, result.Error
	}
	bytePassword := []byte(password)
	byteHashedPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return user, token, err
	}

	tokenDetails, err := authModel.CreateToken(user.ID)
	if err != nil {
		return user, token, err
	}

	token = tokenDetails.AccessToken

	return user, token, nil
}
