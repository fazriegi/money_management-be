package libs

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fazriegi/money_management-be/model"
	"github.com/spf13/viper"
)

type JWT struct {
	secretKey string
	expHour   uint16
}

func InitJWT(viper *viper.Viper) *JWT {
	secretKey := viper.GetString("jwt.key")
	expHour := viper.GetUint16("jwt.expHour")

	return &JWT{
		secretKey,
		expHour,
	}
}

func (s JWT) GenerateJWTToken(user *model.User) (string, error) {
	exp := time.Duration(s.expHour) * time.Hour
	claims := jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s JWT) VerifyJWTTOken(tokenString string) (any, error) {
	errResponse := errors.New("invalid or expired token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}

		return []byte(s.secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, errResponse
	}

	return token.Claims.(jwt.MapClaims), nil
}
