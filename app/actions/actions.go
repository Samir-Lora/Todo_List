package actions

import (
	"todo_list/app/render"
	// modificado
)

var (
	// r is a buffalo/render Engine that will be used by actions
	// on this package to render render HTML or any other formats.
	r = render.Engine
)
