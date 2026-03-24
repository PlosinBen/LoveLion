package handlers

import (
	"net/http"

	"lovelion/internal/utils/errorx"

	"github.com/gin-gonic/gin"
)

// respondError maps an AppError to the appropriate HTTP status and JSON response.
func respondError(c *gin.Context, err error) {
	status := http.StatusInternalServerError

	if errorx.Is(err, errorx.ErrNotFound) {
		status = http.StatusNotFound
	} else if errorx.Is(err, errorx.ErrExpired) || errorx.Is(err, errorx.ErrExhausted) {
		status = http.StatusGone
	} else if errorx.Is(err, errorx.ErrForbidden) {
		status = http.StatusForbidden
	} else if errorx.Is(err, errorx.ErrConflict) {
		status = http.StatusConflict
	} else if errorx.Is(err, errorx.ErrBadRequest) {
		status = http.StatusBadRequest
	}

	c.JSON(status, gin.H{"error": err.Error()})
}
