package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"fucking-fat/internal/handlers/auth/jwt"
	"fucking-fat/internal/models"
	"golang.org/x/crypto/bcrypt"
	"io"

	"net/http"
)

func GetBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetResult(w http.ResponseWriter, res interface{}, status int) {
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		fmt.Printf(fmt.Sprintf("%v", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func DecodeHashPassword(hash, pass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return false, err
	}
	return true, nil
}

func FindToken(r *http.Request) (*models.User, error) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "Token" {
			user, err := jwt.ParseToken(cookie.Value)
			if err != nil {
				if errors.As(err, &jwt.TokenExp) {
					return nil, err
				}
			}
			return user, nil
		}
	}
	return nil, nil
}
