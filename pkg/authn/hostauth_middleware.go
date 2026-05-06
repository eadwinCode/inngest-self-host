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

// SigningKeyOrHostAuthMiddleware allows a request through if it passes either
// the signing key middleware (bearer token from SDKs) or has a valid host auth
// cookie (browser). This is used for endpoints that both the UI and SDKs need.
func SigningKeyOrHostAuthMiddleware(signingKeyMiddleware func(http.Handler) http.Handler, config *HostAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		signingKeyHandler := signingKeyMiddleware(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try host auth cookie first
			if config != nil && config.IsEnabled() {
				if cookie, err := r.Cookie(hostAuthCookieName); err == nil {
					if _, err := config.ValidateToken(cookie.Value); err == nil {
						next.ServeHTTP(w, r)
						return
					}
				}
			}
			// Fall back to signing key auth
			signingKeyHandler.ServeHTTP(w, r)
		})
	}
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
