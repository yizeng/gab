package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/cors"
)

func ConfigCORS(environment string, allowedDomains []string) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowOriginFunc:  createAllowedOriginFunc(allowedDomains),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            strings.EqualFold(environment, "development"),
	})
}

func createAllowedOriginFunc(allowedDomains []string) func(r *http.Request, origin string) bool {
	return func(r *http.Request, origin string) bool {
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
