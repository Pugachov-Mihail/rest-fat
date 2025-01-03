package jwt

import (
	"errors"
	"fmt"
	"fucking-fat/internal/models"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

var TokenExp = errors.New("token expired")

func NewToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()
	claims["auth"] = "service-auth"

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	user := models.User{}
	user.Username = token.Claims.(jwt.MapClaims)["username"].(string)
	user.Id = int64(token.Claims.(jwt.MapClaims)["uid"].(float64))
	exp := int64(token.Claims.(jwt.MapClaims)["exp"].(float64))

	if ValidateTime(exp) {
		return &user, nil
	}

	return nil, TokenExp
}

func ValidateTime(exp int64) bool {
	if time.Unix(exp, 0).UTC().Format("2006-01-02 15:04:05") > time.Now().Format("2006-01-02 15:04:05") {
		return true
	}
	return false
}

func CookiesOnToken(token string, w http.ResponseWriter) {
	cookie := &http.Cookie{Name: "Token", Value: token, Secure: true, Expires: time.Now().Add(time.Hour * 3)}
	http.SetCookie(w, cookie)
}
