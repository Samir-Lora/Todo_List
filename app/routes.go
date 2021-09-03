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
	root.Use(middleware.SetCurrentUser)
	root.Use(middleware.Authorize)
	root.Use(middleware.Authorizeusers)

	root.GET("/", actions.AuthLogin)
	root.POST("/auth", actions.AuthCreate)
	root.GET("/tasks", actions.TaskList)
	root.GET("/task/new", actions.Newtask)
	root.POST("/task/create", actions.Createtask)

	root.GET("/task/newtask", actions.Newtaskuser)
	root.POST("/task/createuser", actions.Createtaskuser)

	root.GET("/task/showtask/{task_id}", actions.Showtask)
	root.GET("/task/edit/{task_id}", middleware.EditTaskAcess(actions.Edittask))
	root.PUT("/task/edit/{task_id}/update", actions.Updatetask)
	root.DELETE("/task/delete/{task_id}", actions.Delete)
	root.PUT("/task/updatecomplete/{task_id}", actions.Updatecomplete)

	root.GET("/users", actions.UserList)
	root.GET("/user/new", actions.NewUser)
	root.POST("/user/create", actions.CreateUser)

	root.GET("/user/new/invitation", actions.NewInvitation)
	root.POST("/user/create/invitation", actions.CreateInvitation)

	root.GET("/user/edit/{user_id}", actions.Edituser)
	root.PUT("/user/edit/{user_id}/update", actions.Updateuser)
	root.DELETE("/user/delete/{user_id}", actions.Deleteuser)
	root.PUT("/user/updateactive/{user_id}", actions.Updateactive)
	root.DELETE("/logout", actions.AuthDestroy)

	root.Middleware.Skip(middleware.Authorize, actions.AuthLogin, actions.AuthCreate, actions.AuthDestroy, actions.NewUser, actions.CreateUser)
	root.Middleware.Skip(middleware.SetCurrentUser, actions.AuthLogin, actions.AuthCreate, actions.AuthDestroy, actions.NewUser)
	root.Middleware.Skip(middleware.Authorizeusers, actions.AuthLogin, actions.AuthCreate, actions.NewUser, actions.TaskList, actions.AuthDestroy, actions.CreateUser, actions.Newtaskuser, actions.Createtask, actions.Createtaskuser, actions.Updatecomplete)
	root.ServeFiles("/", base.Assets)
}
