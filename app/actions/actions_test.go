package actions_test

import (
	"testing"
	"todo_list/app"

	"github.com/gobuffalo/suite/v3"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	bapp := app.New()

	as := &ActionSuite{suite.NewAction(bapp)}
	suite.Run(t, as)
}
