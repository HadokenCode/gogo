package controllers

import (
	"github.com/dolab/gogo"
)

var (
	GettingStart *_GettingStart
)

type _GettingStart struct{}

// @route GET /@gogo/ping
func (_ *_GettingStart) Pong(ctx *gogo.Context) {
	ctx.Logger.Warnf("Visiting domain is: %s", Config.Domain)

	ctx.Text(Config.GettingStart.Greeting)
}
