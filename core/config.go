package core

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

// Config represents migration configuration
type Config struct {
	RootDir  string
	Email    string
	Name     string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewConfig creates a new migration config
func NewConfig() Config {
	root, _ := os.Getwd()
	userEmail, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		panic(err)
	}

	userName, err := exec.Command("git", "config", "--global", "user.name").Output()
	if err != nil {
		panic(err)
	}

	return Config{
		RootDir:  path.Join(root, "migrations"),
		Email:    strings.TrimSpace(string(userEmail)),
		Name:     strings.TrimSpace(string(userName)),
		Host:     "localhost",
		Port:     5432,
		Database: "test",
		Password: "",
		User:     "postgres",
	}
}
