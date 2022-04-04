package models

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	uuid "github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken string
	AccessUUID  string
	AtExpires   int64
}

type AccessDetails struct {
	AccessUUID string
	UserID     int32
}

type AuthModel struct{}

func (m AuthModel) CreateToken(userID int32) (*TokenDetails, error) {

	tokenDetail := &TokenDetails{}
	tokenDetail.AtExpires = time.Now().Add(time.Hour * 72).Unix()
	tokenDetail.AccessUUID = uuid.NewV4().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDetail.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = tokenDetail.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenDetail.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return tokenDetail, nil
}

func (m AuthModel) ExtractToken(request *http.Request) string {
	bearToken := request.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (m AuthModel) VerifyToken(request *http.Request) (*jwt.Token, error) {
	tokenString := m.ExtractToken(request)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m AuthModel) ExtractTokenMetadata(request *http.Request) (*AccessDetails, error) {
	token, err := m.VerifyToken(request)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 32)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     int32(userID),
		}, nil
	}
	return nil, err
}
