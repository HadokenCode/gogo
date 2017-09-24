package controllers

import (
	"testing"
)

func Test_GettingStart_Pong(t *testing.T) {
	testClient.Get(t, "/@gogo/pong")

	testClient.AssertOK()
	testClient.AssertContains(Config.GettingStart.Greeting)
}
