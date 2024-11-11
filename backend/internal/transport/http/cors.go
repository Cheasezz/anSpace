package httpHandlers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func corsMiddleware(origins []string) gin.HandlerFunc {
	corsMdlwr := cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-type"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	})
	
	return corsMdlwr
}
