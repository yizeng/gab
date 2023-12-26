package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ConfigCORS(host string) gin.HandlerFunc {
	conf := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://0.0.0.0") {
				return true
			}

			return strings.HasPrefix(origin, host)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(conf)
}
