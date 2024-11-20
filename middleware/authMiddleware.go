package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	// "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// ContextKey is a type for storing context keys
type ContextKey string

const (
	UserClaimsContextKey ContextKey = "userId"
)

// Middleware holds dependencies for middleware functions
type Middleware struct {
	Log    *zap.Logger
	Config util.Configuration
}

// NewMiddleware creates a new Middleware instance
func NewMiddleware(log *zap.Logger, config util.Configuration) *Middleware {
	return &Middleware{
		Log:    log,
		Config: config,
	}
}

// AuthMiddleware ensures requests are authenticated
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.extractToken(r)
		if err != nil {
			m.Log.Info("Unauthorized access", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
			m.respondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
			return
		}

		claims, err := util.VerifyToken(token, m.Config)
		if err != nil {
			m.Log.Error("Token verification failed", zap.Error(err))
			m.respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add claims to the context
		ctx := context.WithValue(r.Context(), UserClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken retrieves the token from a cookie or the Authorization header
func (m *Middleware) extractToken(r *http.Request) (string, error) {
	// Check for token in cookies
	cookie, err := r.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	// Check for token in Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	return "", http.ErrNoCookie // Reuse standard error for missing token
}

// respondWithError sends an error response with JSON formatting
func (m *Middleware) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
