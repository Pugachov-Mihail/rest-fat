package conf

import (
	"fucking-fat/internal/source"
	"log/slog"
)

type Conf struct {
	Env      string `yaml:"env" default:"dev"`
	Logger   *slog.Logger
	ConfigDB `yaml:"config_db"`
}

type ConfigDB struct {
	Source string `yaml:"source"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Dbname string `yaml:"dbname"`
}

func NewConf() *Conf {
	return &Conf{}
}

func (conf *Conf) DbConf() *source.Posgresql {
	switch {
	case conf.Env == "dev":
		return source.CreateConn(conf.ConfigDB.Source, "sqlite3")
	case conf.Env == "prod":
		return source.CreateConn(conf.ConfigDB.Source, "postgres")
	}
	return nil
}
