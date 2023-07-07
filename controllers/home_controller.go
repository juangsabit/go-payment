package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	// response json
	c.JSON(http.StatusOK, gin.H{
		"name": "Juang Sabit",
		"bio":  "do your best",
	})
}
