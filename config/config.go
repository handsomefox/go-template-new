package config

import (
	"os"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Specification struct {
	Database Database `koanf:"database"`
	Server   Server   `koanf:"server"`
}

type Database struct {
	Host     string `koanf:"host"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	SSLMode  string `koanf:"sslmode"`
	Database string `koanf:"database"`
	Port     int    `koanf:"port"`
}

type Server struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

func New() (*Specification, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config/config.toml"
	}
	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		return nil, err
	}
	var s Specification
	err := k.Unmarshal("", &s)
	return &s, err
}
