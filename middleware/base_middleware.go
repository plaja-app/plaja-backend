package middleware

import "github.com/plaja-app/back-end/config"

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
