package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
)

// Loading connections from database.yml in the pop.Connections

// DB returns the DB connection for the current environment.
func DB() *pop.Connection {
	c, err := pop.Connect(envy.Get("GO_ENV", "development"))
	if err != nil {
		log.Fatal(err)
	}

	return c
}
