package authn

import (
	"net/http"
	"net/url"
)

// SameHostOriginFunc only allows origins that match the request's Host header.
// This prevents cross-origin credentialed requests from arbitrary websites.
func SameHostOriginFunc(r *http.Request, origin string) bool {
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	return parsed.Host == r.Host
}

func HostAuthMiddleware(config *HostAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if config == nil || !config.IsEnabled() {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(hostAuthCookieName)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
				return
			}
			if _, err := config.ValidateToken(cookie.Value); err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
