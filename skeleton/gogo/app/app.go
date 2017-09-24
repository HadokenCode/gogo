package app

import (
	"github.com/skeleton/app/controllers"
	"github.com/skeleton/app/middlewares"
)

type Application struct {
	*controllers.Application

	_ bool
}

func New(runMode, srcPath string) *Application {
	app := &Application{
		Application: controllers.New(runMode, srcPath),
	}

	return app
}

// Middlerwares overwrites gogo.Middleware() method for custom middlewares
// NOTE: DO NOT change the method name, its required by gogo internal!
func (app *Application) Middlewares() {
	// apply your middlewares

	// panic recovery
	app.Use(middlewares.Recovery())
}

// Run overwrites gogo.Run() method by registering middlewares and resources.
// NOTE: DO NOT change the method name, its required by gogo internal!
func (app *Application) Run() {
	// register middlewares
	app.Middlewares()

	// register resources
	app.Resources()

	// run server
	app.Application.Run()
}
