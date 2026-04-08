package app

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
)

func ensureDevAdmin(cfg *config.Config, store repository.Store) error {
	if !cfg.AdminBootstrapEnabled {
		return nil
	}
	if strings.EqualFold(cfg.AppEnv, "production") {
		return nil
	}
	email := strings.TrimSpace(strings.ToLower(cfg.AdminBootstrapEmail))
	password := cfg.AdminBootstrapPassword
	if email == "" || strings.TrimSpace(password) == "" {
		log.Printf("app init: admin bootstrap enabled but credentials are incomplete, skipping")
		return nil
	}

	_, err := store.FindUserByEmail(email)
	if err == nil {
		log.Printf("app init: bootstrap admin already exists: %s", email)
		return nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if _, _, err := store.CreateUser(email, string(hash), domain.RoleAdmin); err != nil {
		return err
	}
	log.Printf("app init: bootstrap admin created: %s", email)
	return nil
}
