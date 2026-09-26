package handler

import (
	"errors"
	"net/http"

	"otomedia/task-managment/internal/apperror"

	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

// RespondError menulis body error JSON yang konsisten di seluruh API:
//
//	{ "error": { "code": "NOT_FOUND", "message": "task not found" } }
//
// Untuk error yang bukan *apperror.AppError (misal error driver DB yang lolos
// tanpa dibungkus), handler tetap membalas 500 INTERNAL_ERROR daripada
// membocorkan detail internal ke client.
func RespondError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, errorResponse{Error: errorBody{Code: appErr.Code, Message: appErr.Message}})
		return
	}
	c.JSON(http.StatusInternalServerError, errorResponse{
		Error: errorBody{Code: "INTERNAL_ERROR", Message: "an unexpected error occurred"},
	})
}

func RespondValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, errorResponse{Error: errorBody{Code: "BAD_REQUEST", Message: message}})
}
