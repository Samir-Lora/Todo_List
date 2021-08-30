package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func Login(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("users/newuser.plush.html"))

}
