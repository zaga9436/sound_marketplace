package service

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/soundmarket/backend/internal/auth"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
)

type AuthService struct {
	store repository.Store
	jwt   *auth.JWTManager
}

type AuthResult struct {
	Token   string         `json:"token,omitempty"`
	User    domain.User    `json:"user"`
	Profile domain.Profile `json:"profile"`
}

func NewAuthService(store repository.Store, jwt *auth.JWTManager) *AuthService {
	return &AuthService{store: store, jwt: jwt}
}

func (s *AuthService) Register(email, password string, role domain.Role) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}
	if role != domain.RoleCustomer && role != domain.RoleEngineer {
		return nil, fmt.Errorf("only customer or engineer can self-register")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, profile, err := s.store.CreateUser(email, string(hash), role)
	if err != nil {
		return nil, err
	}
	token, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, User: user, Profile: profile}, nil
}

func (s *AuthService) Login(email, password string) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("email and password are required")
	}
	user, err := s.store.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	profile, err := s.store.GetProfile(user.ID)
	if err != nil {
		return nil, err
	}
	token, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, User: user, Profile: profile}, nil
}

func (s *AuthService) Me(userID string) (*AuthResult, error) {
	user, err := s.store.GetUser(userID)
	if err != nil {
		return nil, err
	}
	profile, err := s.store.GetProfile(userID)
	if err != nil {
		return nil, err
	}
	return &AuthResult{User: user, Profile: profile}, nil
}
