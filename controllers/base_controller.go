package controllers

import "github.com/plaja-app/back-end/config"

// Controller the repository used by the controllers.
var Controller *BaseController

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

// NewControllers sets the repository for the controllers.
func NewControllers(c *BaseController) {
	Controller = c
}
