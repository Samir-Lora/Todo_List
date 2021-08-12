package app

import (
	base "todo_list"
	"todo_list/app/actions"
	"todo_list/app/actions/home"
	"todo_list/app/middleware"

	"github.com/gobuffalo/buffalo"
)

// SetRoutes for the application
func setRoutes(root *buffalo.App) {
	root.Use(middleware.Transaction)
	root.Use(middleware.ParameterLogger)
	root.Use(middleware.CSRF)

	root.GET("/", home.Index)
	root.GET("/task", actions.TaskList)
	root.GET("/task/new", actions.Newtask)
	root.POST("/task/create", actions.Createtask)
	root.GET("/task/showtask/{task_ID}", actions.Showtask)
	root.GET("/task/edit/{task_ID}", actions.Edittask)
	root.PUT("/task/edit/{task_ID}/update", actions.Updatetask)
	root.DELETE("/task/delete/{task_ID}", actions.Delete)
	root.PUT("/task/updatecomplete/{task_ID}", actions.Updatecomplete)
	root.ServeFiles("/", base.Assets)
}
