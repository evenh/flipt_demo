package actions

import (
	"github.com/gobuffalo/buffalo"
)

func HomeHandler(c buffalo.Context) error {
	return c.Redirect(302, "/items")
}
