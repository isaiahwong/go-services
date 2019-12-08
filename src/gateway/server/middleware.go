package server

import (
	"github.com/gin-gonic/gin"
)

func notFound(c *gin.Context) {
	c.JSON(404, gin.H{
		"success": false,
		"error":   "not_found",
	})
}
