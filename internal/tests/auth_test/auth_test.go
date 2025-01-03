package auth_test

import (
	"fucking-fat/internal/models"
	"fucking-fat/internal/source/mocks"
	"go.uber.org/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	//client := http.Client{Transport: http.DefaultTransport}
	crl := gomock.NewController(t)
	dbs := mocks.NewMockServiceDbs(crl)
	user := &models.User{Username: "admin", Password: "1234"}
	user2 := &models.User{Id: 1, Username: "admin", Password: "1234"}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	dbs.EXPECT().Register(user, log).
		Return(user2, nil)

}
