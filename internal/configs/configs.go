package configs

import (
	"fmt"
	"github.com/ddritzenhoff/dindin/internal/http/rest"
	"log"
	"os"
	"strings"
)

// Configs handles all dependencies required for handling configurations
type Configs struct{}

func (cfg *Configs) DBName() (string, error) {
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Print("DB_NAME is not set")
		return "", fmt.Errorf("DB_NAME is not set")
	}
	return strings.TrimSpace(dbName), nil
}

func (cfg *Configs) HTTP() (*rest.Config, error) {
	return &rest.Config{
		Host: "localhost",
		Port: "8080",
	}, nil
}

func (cfg *Configs) SlackConfig() (*SlackConfig, error) {
	return NewSlackConfig()
}

func NewConfigService() (*Configs, error) {
	return &Configs{}, nil
}
