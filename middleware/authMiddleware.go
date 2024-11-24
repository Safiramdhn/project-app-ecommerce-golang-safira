package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
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
			m.handleUnauthorized(w, r, "Unauthorized access: "+err.Error())
			return
		}

		// Log successful extraction of token
		m.Log.Info("Token extracted successfully",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("token", token),
		)

		claims, err := util.VerifyToken(token, m.Config)
		if err != nil {
			m.handleUnauthorized(w, r, "Invalid token: "+err.Error())
			return
		}

		// Log claims parsed successfully
		m.Log.Info("Claims parsed successfully",
			zap.String("userId", claims["userId"].(string)),
			zap.Time("issuedAt", claims["iat"].(time.Time)),
			zap.Time("expiresAt", claims["exp"].(time.Time)),
		)

		user := model.User{
			ID: claims["userId"].(string),
		}

		// Add user info to the context and log context update
		ctx := context.WithValue(r.Context(), UserClaimsContextKey, user)
		m.Log.Info("User added to context", zap.String("userId", user.ID))

		// Proceed with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken retrieves the token from a cookie or the Authorization header
func (m *Middleware) extractToken(r *http.Request) (string, error) {
	// Check for token in Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	// Log missing token
	m.Log.Warn("Missing token in cookie or Authorization header",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	)
	return "", errors.New("missing token in cookie or Authorization header")
}

// handleUnauthorized handles unauthorized access with logging and response
func (m *Middleware) handleUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	// Log unauthorized access attempt
	m.Log.Warn("Unauthorized access attempt",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("message", message),
	)

	// Send error response
	m.respondWithError(w, http.StatusUnauthorized, message)
}

// respondWithError sends an error response with JSON formatting
func (m *Middleware) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	// Log the error response
	m.Log.Error("Responding with error",
		zap.Int("statusCode", statusCode),
		zap.String("message", message),
	)

	// Prepare and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
