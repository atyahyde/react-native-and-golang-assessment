package apperror

import "net/http"

// AppError adalah error terstruktur yang dipetakan langsung ke response HTTP,
// supaya semua error di API punya bentuk yang konsisten:
//
//	{ "error": { "code": "...", "message": "..." } }
type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict dipakai untuk kasus judul task duplikat (Task 4: harus 409, bukan 500).
func Conflict(message string) *AppError {
	return New(http.StatusConflict, "CONFLICT", message)
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
