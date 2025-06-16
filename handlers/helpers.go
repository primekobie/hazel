package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getUUIDparam(c *gin.Context, key string) (uuid.UUID, error) {
	idString := c.Param(key)
	return uuid.Parse(idString)
}
