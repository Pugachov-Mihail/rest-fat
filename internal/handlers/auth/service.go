package auth

import (
	"errors"
	"fmt"
	"fucking-fat/internal/handlers/auth/jwt"
	"fucking-fat/internal/helpers"
	"fucking-fat/internal/models"
	"fucking-fat/internal/source"
	"log/slog"
)

type ServiceAuth struct {
	log *slog.Logger
	db  source.Posgresql
}

func (s *ServiceAuth) Login(username, password string) (string, error) {
	log := s.log.With("login", username)

	user, err := s.db.Login(username, password, log)
	if err != nil {
		log.Warn("Error equal profile")
		return "", err
	}

	if flag, err := helpers.DecodeHashPassword(user.Password, password); !flag {
		log.Warn("Invalid password", err)
	}

	token, err := jwt.NewToken(user)
	if err != nil {
		log.Warn("Error new token", err)
		return "", err
	}
	return token, nil
}

func (s *ServiceAuth) Logout(token string) error {
	log := s.log.With("service", "logout")
	user, err := jwt.ParseToken(token)
	if err != nil {
		if errors.Is(err, jwt.TokenExp) {
			log.Warn(fmt.Sprintf("%s:%s", jwt.TokenExp, err))
			return jwt.TokenExp
		}
		log.Warn("Error parse token", err)
	}
	//TODO После добавление кеша хранить токены там, поиск токена будет по username
	_, err = s.db.Logout(user.Username, log)
	if err != nil {
		log.Warn("Error logout", err)
	}
	return nil
}

func (s *ServiceAuth) Register(username, password, email string) (string, error) {
	log := s.log
	pass, err := helpers.HashPassword(password)
	if err != nil {
		log.Warn("Failed to hash password")
	}
	var user = models.User{Username: username, Password: pass, Email: email}
	data, err := s.db.Register(&user, log)

	if err != nil {
		log.Warn("Error register", err)
	}

	token, err := jwt.NewToken(data)
	if err != nil {
		log.Warn("Error new token", err)
	}
	log.Info("User created", user.Username)
	return token, nil
}

func (s *ServiceAuth) SaveUserInfo(data RequestUserData) error {
	log := s.log
	user, err := s.db.UserInfo(
		models.User{
			Id:        data.Id,
			Username:  data.Username,
			FirstName: data.FirstName,
			LastName:  data.LastName,
		},
		log)
	if err != nil {
		log.Warn("Error get profile info", err)
	}
	if user == nil {
		log.Warn("Error create profile info")
		return errors.New("error create profile info")
	}
	return nil
}
