package middleware

import (
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ConfigCORS(allowedDomains []string) gin.HandlerFunc {
	conf := cors.Config{
		AllowOriginFunc:  createAllowedOriginFunc(allowedDomains),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(conf)
}

func createAllowedOriginFunc(allowedDomains []string) func(origin string) bool {
	return func(origin string) bool {
		o, err := url.Parse(origin)
		if err != nil {
			return false
		}
		hostname := o.Hostname()

		localDomains := []string{"localhost", "127.0.0.1", "0.0.0.0"}
		allowedDomains = append(allowedDomains, localDomains...)
		for _, domain := range allowedDomains {
			if strings.EqualFold(hostname, domain) || strings.HasSuffix(hostname, "."+domain) {
				return true
			}
		}

		return false
	}
}
