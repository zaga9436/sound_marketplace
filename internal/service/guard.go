package service

import (
	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
)

func ensureActiveUser(store repository.Store, actor domain.User) error {
	if actor.Role == domain.RoleAdmin {
		return nil
	}
	user, err := store.GetUser(actor.ID)
	if err != nil {
		return apierr.NotFound("user not found")
	}
	if user.IsSuspended {
		return apierr.Forbidden("user is suspended")
	}
	return nil
}
