package authn

import (
	"encoding/json"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authStatusResponse struct {
	AuthRequired  bool `json:"authRequired"`
	Authenticated bool `json:"authenticated"`
}

func HostAuthLoginHandler(config *HostAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config == nil || !config.IsEnabled() {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "authentication not enabled"})
			return
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if !config.ValidateCredentials(req.Email, req.Password) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}

		token, err := config.CreateToken(req.Email)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create session"})
			return
		}

		http.SetCookie(w, sessionCookie(r, token, int(hostAuthExpiry.Seconds())))

		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func HostAuthLogoutHandler(config *HostAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, sessionCookie(r, "", -1))
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func HostAuthStatusHandler(config *HostAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := authStatusResponse{
			AuthRequired:  config != nil && config.IsEnabled(),
			Authenticated: false,
		}

		if resp.AuthRequired {
			if cookie, err := r.Cookie(hostAuthCookieName); err == nil {
				if _, err := config.ValidateToken(cookie.Value); err == nil {
					resp.Authenticated = true
				}
			}
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

func isSecureRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return r.Header.Get("X-Forwarded-Proto") == "https"
}

func sessionCookie(r *http.Request, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     hostAuthCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
