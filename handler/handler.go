package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type Func[T any] func(w http.ResponseWriter, r *http.Request) (*T, error)

func Create[T any](fn Func[T]) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, err := fn(w, r)
		if err != nil {
			var apierror *Error
			if !errors.As(err, &apierror) {
				apierror = NewErrorFromError(err, http.StatusInternalServerError)
			}
			WriteError(r.Context(), w, apierror)
			return
		}

		if err := WriteJSON(w, response, http.StatusOK); err != nil {
			slog.LogAttrs(r.Context(), slog.LevelError, "Failed to write JSON to response", slog.Any("error", err))
		}
	})
}

func Decode[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()

	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, NewErrorFromError(err, http.StatusInternalServerError)
	}

	return &t, nil
}

type Error struct {
	Message string
	Status  int
}

func (e Error) Error() string {
	return e.Message + ", status: " + http.StatusText(e.Status)
}

func NewError(message string, status int) *Error {
	return &Error{
		Message: message,
		Status:  status,
	}
}

func NewErrorFromStatus(status int) *Error {
	return &Error{
		Message: http.StatusText(status),
		Status:  status,
	}
}

func NewErrorFromError(err error, status int) *Error {
	return &Error{
		Message: err.Error(),
		Status:  status,
	}
}

func NewInternalServerError() *Error {
	return &Error{
		Message: http.StatusText(http.StatusInternalServerError),
		Status:  http.StatusInternalServerError,
	}
}

func WriteJSON(w http.ResponseWriter, v any, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(ctx context.Context, w http.ResponseWriter, err *Error) {
	if err := WriteJSON(w, map[string]any{
		"error":  err.Message,
		"status": http.StatusText(err.Status),
	}, err.Status); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "Failed to write JSON error", slog.Any("error", err))
	}
}
