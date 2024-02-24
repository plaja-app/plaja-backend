package middleware

import "github.com/plaja-app/back-end/config"

// Middleware the repository used by the middlware.
var Middleware *BaseMiddleware

// BaseMiddleware holds the base information needed for middleware.
type BaseMiddleware struct {
	App *config.AppConfig
}

// NewBaseMiddleware creates a new BaseMiddleware.
func NewBaseMiddleware(app *config.AppConfig) *BaseMiddleware {
	return &BaseMiddleware{
		App: app,
	}
}

// NewMiddleware sets the repository for the middleware.
func NewMiddleware(m *BaseMiddleware) {
	Middleware = m
}
