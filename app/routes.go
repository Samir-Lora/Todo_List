package app

import (
	base "todo_list"
	"todo_list/app/actions"
	"todo_list/app/middleware"

	"github.com/gobuffalo/buffalo"
)

// SetRoutes for the application
func setRoutes(root *buffalo.App) {
	root.Use(middleware.Transaction)
	root.Use(middleware.ParameterLogger)
	root.Use(middleware.CSRF)
	root.Use(middleware.IncompleteTask)
	root.Use(middleware.Datenow)

	root.GET("/", actions.TaskList)
	root.GET("/task/new", actions.Newtask)
	root.POST("/task/create", actions.Createtask)
	root.GET("/task/showtask/{task_id}", actions.Showtask)
	root.GET("/task/edit/{task_id}", middleware.EditTaskAcess(actions.Edittask))
	root.PUT("/task/edit/{task_id}/update", actions.Updatetask)
	root.DELETE("/task/delete/{task_id}", actions.Delete)
	root.PUT("/task/updatecomplete/{task_id}", actions.Updatecomplete)
	root.ServeFiles("/", base.Assets)
}
