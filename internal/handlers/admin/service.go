package admin

import (
	"errors"
	"fmt"
	"fucking-fat/internal/handlers/auth/jwt"
	"fucking-fat/internal/helpers"
	"fucking-fat/internal/models"
	"fucking-fat/internal/source"
	"log/slog"
)

type ServiceAdmin struct {
	Log *slog.Logger
	db  source.Posgresql
}

type AdminDb interface {
	GetAllUsersAdmin() ([]models.User, error)
}

func (a *ServiceAdmin) Login(username, password string) (string, error) {
	log := a.Log.With("service-login", "admin")
	user, err := a.db.Login(username, password, log)
	if err != nil {
		log.Warn("Error while truing profile", err)
		return "", err
	}
	flag, _ := helpers.DecodeHashPassword(user.Password, password)
	if !flag {
		log.Warn("Error while truing profile")
		return "", errors.New("invalid username or password")
	}
	perm := source.MapPermissions(user.Role)
	if perm == models.PermissionAdmin {
		token, err := jwt.NewToken(user)
		if err != nil {
			log.Warn("Error generate token to admin", err)
			return "", err
		}
		return token, nil
	} else {
		log.Warn("Error permission profile", user.Username)
		return "", err
	}
}

func (a *ServiceAdmin) Logout(token string) error {
	log := a.Log.With("service", "logout")
	user, err := jwt.ParseToken(token)
	if err != nil {
		if errors.Is(err, jwt.TokenExp) {
			log.Warn(fmt.Sprintf("%s:%s", jwt.TokenExp, err))
			return jwt.TokenExp
		}
		log.Warn("Error parse token", err)
	}
	//TODO После добавление кеша хранить токены там, поиск токена будет по username
	_, err = a.db.Logout(user.Username, log)
	if err != nil {
		log.Warn("Error logout", err)
	}
	return nil
}

func (a *ServiceAdmin) AllUsersIntoRole() ([]models.User, error) {
	log := a.Log
	users, err := a.db.GetAllUsersAdmin()
	if err != nil {
		log.Warn("Error getting all users", err)
		return nil, err
	}
	return users, nil
}
