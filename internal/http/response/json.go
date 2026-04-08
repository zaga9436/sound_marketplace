package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/repository"
)

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func FromError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apierr.ErrUnauthorized):
		Error(w, http.StatusUnauthorized, apierr.Message(err))
	case errors.Is(err, apierr.ErrForbidden):
		Error(w, http.StatusForbidden, apierr.Message(err))
	case errors.Is(err, apierr.ErrNotFound), errors.Is(err, repository.ErrNotFound):
		Error(w, http.StatusNotFound, apierr.Message(err))
	case errors.Is(err, apierr.ErrConflict):
		Error(w, http.StatusConflict, apierr.Message(err))
	case errors.Is(err, apierr.ErrBadRequest):
		Error(w, http.StatusBadRequest, apierr.Message(err))
	default:
		Error(w, http.StatusInternalServerError, err.Error())
	}
}
