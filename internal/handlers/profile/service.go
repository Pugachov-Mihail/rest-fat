package profile

import (
	"fucking-fat/internal/models"
	"fucking-fat/internal/source"
	"log/slog"
)

type ServiceUser struct {
	log *slog.Logger
	db  source.Posgresql
}

type Profile interface {
	ProfileInfo(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)

}

func (s *ServiceUser) Profile(user models.User)  {
	log := s.log

	user, err := s.db.
}
