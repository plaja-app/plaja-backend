package controllers

import "github.com/plaja-app/back-end/config"

// BaseController holds the base information needed for all the controllers.
type BaseController struct {
	App *config.AppConfig
}

// NewBaseController creates a new BaseController.
func NewBaseController(app *config.AppConfig) *BaseController {
	return &BaseController{
		App: app,
	}
}
