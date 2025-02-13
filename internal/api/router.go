package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/hashocr/internal/api/handlers"
)

func NewRouter(handler *handlers.Handler) http.Handler {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.GET("/health", handler.Ping)
	v1.POST("/pdf2json", handler.CreateJsonFromPdf)

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "Method Not Allowed",
		})
	})

	return router
}
