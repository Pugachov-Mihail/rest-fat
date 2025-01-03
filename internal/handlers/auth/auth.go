package auth

import (
	"encoding/json"
	"fmt"
	"fucking-fat/internal/handlers/auth/jwt"
	"fucking-fat/internal/helpers"
	"log/slog"
	"net/http"
	"strings"
)

type HTTPAuth struct {
	Name     string
	Log      *slog.Logger
	Services ServiceAuth
}

func NewHttpAuth(name string, log *slog.Logger) *HTTPAuth {
	return &HTTPAuth{Name: name, Log: log, Services: ServiceAuth{log: log}}
}

func (h *HTTPAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	auth := fmt.Sprintf("/%s", strings.ToLower(h.Name))
	switch {
	case path == auth+"/login":
		h.Auth(w, r)
	case path == auth+"/update-profile":
		h.UserInfo(w, r)
	case path == auth+"/logout":
		h.Logout(w, r)
	case path == auth+"/register":
		h.Register(w, r)
	default:
		h.Page404(w, r)
	}
}

func (h *HTTPAuth) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log := h.Log.With("service", "auth")

		body, err := helpers.GetBody(r)
		if err != nil {
			log.Error(fmt.Sprintf("Error getting body: %s", err))
		}

		var data RequestLogin
		err = json.Unmarshal(body, &data)
		res := data.Validate()
		if err != nil {
			log.Error(fmt.Sprintf("Error umrashal body: %s", err))
		}

		if res != "" {
			log.Warn(fmt.Sprintf("Error validate body: %s", res))
			helpers.GetResult(w, res, http.StatusBadRequest)
			return
		}

		token, err := h.Services.Login(data.Username, data.Pass)
		if err != nil {
			helpers.GetResult(w, struct {
				message string
			}{fmt.Sprintf("%s", err)}, http.StatusForbidden)
			log.Error(fmt.Sprintf("Error login: %s", err))
		}

		log.Info(fmt.Sprintf("Successfully auth profile: %s", data.Username))
		jwt.CookiesOnToken(token, w)
		w.WriteHeader(200)

	} else {
		helpers.GetResult(w, "No post method", http.StatusForbidden)
	}
}

func (h *HTTPAuth) UserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log := h.Log.With("service", "profile-info")

		body, err := helpers.GetBody(r)
		if err != nil {
			fmt.Println("Error getting body", err)
		}

		var data RequestUserData
		err = json.Unmarshal(body, &data)
		res := data.Validate()
		if err != nil {
			log.Error(fmt.Sprintf("Error parsing form: %s", err))
		}
		if res != "" {
			log.Warn(fmt.Sprintf("Error validate body: %s", res))
			helpers.GetResult(w, res, http.StatusBadRequest)
			return
		}

		//TODO Отрефакторить, вынести поиск токена в мидлвару
		user, err := helpers.FindToken(r)
		if err != nil {
			log.Error(fmt.Sprintf("Error finding token: %s", err))
			return
		}
		data.Id, data.Username = user.Id, user.Username
		err = h.Services.SaveUserInfo(data)
		if err != nil {
			log.Warn(fmt.Sprintf("Error saving profile: %s", err))
			helpers.GetResult(w, struct{ message string }{"Error saving profile"}, http.StatusForbidden)
			return
		}

		log.Info(fmt.Sprintf("Successfully profile info: %s", data.Username))
		helpers.GetResult(w, map[string]string{
			"firstName": data.FirstName,
			"last_name": data.LastName,
		}, http.StatusCreated)
	} else {
		helpers.GetResult(w, "No post method", http.StatusForbidden)
	}
}

func (h *HTTPAuth) Logout(w http.ResponseWriter, r *http.Request) {
	log := h.Log.With("service", "logout")
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "Token" {
			err := h.Services.Logout(cookie.Value)
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

func (h *HTTPAuth) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log := h.Log.With("service", "register")
		body, err := helpers.GetBody(r)
		if err != nil {
			log.Error(fmt.Sprintf("Error getting body: %s", err))
		}
		var data RequestRegister
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Error(fmt.Sprintf("Error parsing form: %s", err))
		}
		res := data.Validate()
		if res != "" {
			log.Error(fmt.Sprintf("Error validate body: %s", err))
			helpers.GetResult(w, struct{ Message string }{res}, http.StatusBadRequest)
			return
		}
		user, err := h.Services.Register(data.Username, data.Password, data.Email)
		if err != nil {
			log.Error(fmt.Sprintf("Error register: %s", err))
		}
		jwt.CookiesOnToken(user, w)

		helpers.GetResult(w, struct {
			Message string
		}{"User create"}, http.StatusCreated)
	} else {
		helpers.GetResult(w, "No post method", http.StatusForbidden)
	}
}

func (h *HTTPAuth) Page404(w http.ResponseWriter, r *http.Request) {
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
