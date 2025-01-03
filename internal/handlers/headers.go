package handlers

import "net/http"

type Header interface {
	Auth(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Page404(w http.ResponseWriter, r *http.Request)
}

type Service interface {
	Login(username, password string) (string, error)
	Logout(token string) error
}
