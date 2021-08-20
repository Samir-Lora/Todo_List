package app

import (
	"testing"

	"github.com/gobuffalo/buffalo"
)

func Test_app(t *testing.T) {
	app := buffalo.New(buffalo.Options{})
	if app.Env == "" {
		t.Fatalf(app.Env)
	}
	if app.SessionName == "" {
		t.Fatalf(app.SessionName)
	}

}
