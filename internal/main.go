package main

import (
	"fucking-fat/conf"
	"fucking-fat/internal/handlers/admin"
	"fucking-fat/internal/handlers/auth"
	"fucking-fat/internal/handlers/profile"
	"log/slog"
	"net/http"
	"os"
)

func middlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	cfg := conf.NewConf()
	cfg.DbConf()

	cfg.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log := cfg.Logger.With("Main server", "fucking-fat")

	authHttp := auth.NewHttpAuth("Auth", cfg.Logger)
	adminHttp := admin.NewHttpAdmin("Admin", cfg.Logger)
	userHttp := profile.NewHttpUser("Profile", cfg.Logger)

	mux.Handle("/auth/", middlewareJson(authHttp))
	mux.Handle("/admin/", middlewareJson(adminHttp))
	mux.Handle("/profile/", middlewareJson(userHttp))

	//http.HandleFunc("/profile", auth.UserInfo)
	log.Info("Starting server on port 80")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Error("Error in ListenAndServe: ", err)
		panic(err)
	}
}
