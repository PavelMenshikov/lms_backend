package httperror

import (
	"log/slog"
	"net/http"
)

func Internal(w http.ResponseWriter, err error) {
	slog.Error("internal server error", slog.String("error", err.Error()))
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, err error) {
	slog.Warn("bad request", slog.String("error", err.Error()))
	http.Error(w, "bad request", http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter, err error) {
	slog.Warn("not found", slog.String("error", err.Error()))
	http.Error(w, "not found", http.StatusNotFound)
}

func Unauthorized(w http.ResponseWriter, err error) {
	slog.Warn("unauthorized", slog.String("error", err.Error()))
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter) {
	http.Error(w, "forbidden", http.StatusForbidden)
}

func Conflict(w http.ResponseWriter, err error) {
	slog.Warn("conflict", slog.String("error", err.Error()))
	http.Error(w, "conflict", http.StatusConflict)
}

func Status(w http.ResponseWriter, err error, status int) {
	slog.Warn("request error", slog.String("error", err.Error()), slog.Int("status", status))
	http.Error(w, http.StatusText(status), status)
}
