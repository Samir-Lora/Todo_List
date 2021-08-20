package models

import (
	"testing"

	"github.com/gobuffalo/suite/v3"
)

type ModelSuite struct {
	suite.Model
}

func Test_Model(t *testing.T) {
	suite.Run(t, &ModelSuite{
		Model: *suite.NewModel(),
	})
}
