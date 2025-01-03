package profile

import (
	"encoding/json"
	"fmt"
	"fucking-fat/internal/helpers"
	"fucking-fat/internal/models"
	"log/slog"
	"net/http"
	"strings"
)

type UserHttp struct {
	Name    string
	log     *slog.Logger
	Service ServiceUser
}

type Message struct {
	Message string `json:"message"`
}

func NewHttpUser(name string, logger *slog.Logger) *UserHttp {
	return &UserHttp{
		Name:    name,
		log:     logger,
		Service: ServiceUser{log: logger},
	}
}

func (u *UserHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log := u.log.With("profile-service", "profile")
	user := fmt.Sprintf("/%s/", strings.ToLower(u.Name))

	switch {
	case path == user:
		u.ProfileUser(w, r)
	}
}

func (u *UserHttp) ProfileUser(w http.ResponseWriter, r *http.Request) {
	log := u.log.With("profile-service", "profile")
	body, err := helpers.GetBody(r)
	if err != nil {
		log.Error("Error reading body", err)
		helpers.GetResult(w, Message{"Error reading body"}, http.StatusBadRequest)
	}
	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Error("Error unmarshalling body", err)
		helpers.GetResult(w, Message{"Error unmarshalling body"}, http.StatusBadRequest)
	}

}
