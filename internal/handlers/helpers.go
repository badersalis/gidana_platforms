package handlers

import (
	"strconv"

	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/gin-gonic/gin"
)

// paramUint parses a named URL param as uint. Returns 0 if invalid.
func paramUint(c *gin.Context, name string) uint {
	v, _ := strconv.ParseUint(c.Param(name), 10, 64)
	return uint(v)
}

// handleErr writes the HTTP error for a ServiceError and returns true, or returns false if err is nil.
func handleErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if se, ok := services.IsServiceError(err); ok {
		c.JSON(se.Code, gin.H{"success": false, "error": se.Message})
		return true
	}
	c.JSON(500, gin.H{"success": false, "error": "Internal server error"})
	return true
}
