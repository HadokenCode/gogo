package controllers

import (
	"github.com/dolab/gogo"
)

type Application struct {
	*gogo.AppServer
}

func New(runMode, srcPath string) *Application {
	appServer := gogo.New(runMode, srcPath)

	err := NewAppConfig(appServer.Config())
	if err != nil {
		panic(err.Error())
	}

	return &Application{appServer}
}

// Resources overwrites gogo.Resources() methods for custom resources.
// NOTE: DO NOT change the method name, its required by gogo!
func (app *Application) Resources() {
	// register your resources
	// app.GET("/", handler)

	app.GET("/@gogo/ping", GettingStart.Pong)
}
