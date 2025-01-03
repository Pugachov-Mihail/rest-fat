package source

import (
	"database/sql"
	"fmt"
	"fucking-fat/internal/models"

	"log/slog"
)

//go:generate mockgen -source=posgresql.go -destination=mocks/mock_database.go -package=mocks

type Posgresql struct {
	Db *sql.DB
}

func CreateConn(db string, driver string) *Posgresql {
	conn, err := sql.Open(driver, db)
	if err != nil {
		panic(err)
	}

	return &Posgresql{conn}
}

func (p *Posgresql) Login(username, password string, log *slog.Logger) (*models.User, error) {
	defer func() {
		if err := p.Db.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	row, err := p.Db.Query("INSERT INTO users(username, password) VALUES(?, ?)", username, password)
	if err != nil {
		log.Error("Error inserting profile", err)
		return nil, err
	}
	var user models.User
	err = row.Scan(&user)
	if err != nil {
		log.Error("Error scanning profile", err)
		return nil, err
	}
	return &models.User{Username: user.Username, Id: user.Id}, nil
}

func (p *Posgresql) Logout(username string, log *slog.Logger) (*models.User, error) {
	defer func() {
		if err := p.Db.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	row, err := p.Db.Query("SELECT id FROM users WHERE username = ?", username)
	if err != nil {
		log.Error("Error selecting profile", err)
		return nil, err
	}
	var user models.User
	err = row.Scan(&user)
	if err != nil {
		log.Error("Error scanning profile", err)
		return nil, err
	}
	return &models.User{Username: user.Username, Id: user.Id}, nil
}

func (p *Posgresql) Register(user *models.User, log *slog.Logger) (*models.User, error) {
	defer func() {
		if err := p.Db.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	row, err := p.Db.Query("INSERT INTO users(username, password) VALUES(?, ?) RETURN id;", user.Username, user.Password)
	if err != nil {
		log.Error("Error inserting profile", err)
		return nil, err
	}
	err = row.Scan(&user)
	if err != nil {
		log.Error("Error scanning profile", err)
		return nil, err
	}

	return user, nil
}

func (p *Posgresql) UserInfo(data models.User, log *slog.Logger) (*models.User, error) {
	defer func() {
		if err := p.Db.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	row, err := p.Db.Query("SELECT id, username FROM users WHERE username = ?", data.Username)
	if err != nil {
		log.Error("Error selecting profile", err)
		return nil, err
	}
	var user models.User
	err = row.Scan(&user)
	if err != nil {
		log.Error("Error scanning profile", err)
		return nil, err
	}
	return &user, nil
}

func (p *Posgresql) GetAllUsersAdmin() ([]models.User, error) {
	panic("implement me")
}
