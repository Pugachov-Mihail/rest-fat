package source

import (
	"fucking-fat/internal/models"
	"log/slog"
)

//go:generate mockgen -source=posgresql.go -destination=./mock_database.go -package=mocks

type ServiceDbs interface {
	Register(user *models.User, log *slog.Logger) (*models.User, error)
	UserInfo(data models.User, log *slog.Logger) (*models.User, error)
}

type AuthDbs interface {
	Login(username, password string, log *slog.Logger) (*models.User, error)
	Logout(username string, log *slog.Logger) (*models.User, error)
}
