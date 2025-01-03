package admin

import (
	"encoding/json"
	"fmt"
	"fucking-fat/internal/handlers/auth/jwt"
	"fucking-fat/internal/helpers"
	"fucking-fat/internal/models"
	"log/slog"
	"net/http"
	"strings"
)

type Admin struct {
	Name    string
	Log     *slog.Logger
	Service ServiceAdmin
}

func NewHttpAdmin(name string, logger *slog.Logger) *Admin {
	return &Admin{name, logger, ServiceAdmin{Log: logger}}
}

func (a *Admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	admin := fmt.Sprintf("/%s", strings.ToLower(a.Name))

	switch {
	case path == admin+"/login":
		a.Auth(w, r)
	case path == admin+"/logout":
		a.Logout(w, r)
	case path == admin+"/get-users":
		a.GetAllUsers(w, r)
	default:
		a.Page404(w, r)
	}
}

func (a *Admin) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log := a.Log.With("service-auth", "admin")
		body, err := helpers.GetBody(r)
		if err != nil {
			log.Error("Error reading body", err)
		}
		var user models.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Error("Error unmarshalling body", err)
			return
		}

		value, err := a.Service.Login(user.Username, user.Password)
		if err != nil {
			log.Error("Error login", err)
			helpers.GetResult(w, struct {
				message string
			}{"Error login"}, http.StatusUnauthorized)
		}
		jwt.CookiesOnToken(value, w)
	} else {
		helpers.GetResult(w, "Not POST method", http.StatusMethodNotAllowed)
	}
}

func (a *Admin) Logout(w http.ResponseWriter, r *http.Request) {
	log := a.Log.With("service-admin", "logout")
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "Token" {
			err := a.Service.Logout(cookie.Value)
			if err != nil {
				log.Error(fmt.Sprintf("Error logout: %s", err))
				helpers.GetResult(w, fmt.Sprintf("%s", err), http.StatusForbidden)
			}
			helpers.GetResult(w, "Logout", http.StatusOK)
		} else {
			helpers.GetResult(w, "No token", http.StatusGone)
		}
	}
}

func (a *Admin) Page404(w http.ResponseWriter, r *http.Request) {
	var res = make(map[string]string)
	res["message"] = "Not found"
	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(404)
	if _, err := w.Write(data); err != nil {
		fmt.Println(err)
	}
}

func (a *Admin) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log := a.Log.With("service: get_all_users", "admin")
	if r.Method == "GET" {
		users, err := a.Service.AllUsersIntoRole()
		if err != nil {
			log.Warn("Error getting all users from database", err)
		}
		helpers.GetResult(w, struct {
			users []models.User
		}{users}, http.StatusOK)
	} else if r.Method == "POST" {
		body, err := helpers.GetBody(r)
		if err != nil {
			log.Error("Error reading body", err)
		}
		var user models.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Error("Error unmarshalling body", err)
		}
		_, err = a.Service.db.Register(&user, log)
		if err != nil {
			log.Error("Error registering profile", err)
		}
		users, err := a.Service.AllUsersIntoRole()
		if err != nil {
			log.Warn("Error getting all users from database", err)
		}
		helpers.GetResult(w, struct {
			users []models.User
		}{users}, http.StatusOK)
	}
}
