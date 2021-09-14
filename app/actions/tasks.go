package actions

import (
	"net/http"
	"time"
	"todo_list/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5" // modificado
	"github.com/gofrs/uuid"
)

func TaskList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	status := c.Param("complete")
	user := c.Value("current_user").(models.User)
	q := tx.PaginateFromParams(c.Params())

	if user.Rol == "user" {
		q.Where("user_id = ?", user.ID)
	}
	if status != "" {
		q.Where("complete = ?", status)
	}

	tasks := models.Tasks{}
	if err := q.Eager("User").Order("priority, task, date").All(&tasks); err != nil {
		return err
	}

	c.Set("taskspaginator", q.Paginator)
	c.Set("user", user)
	c.Set("tasks", tasks)

	return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
}

func Newtask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	users := models.Users{}
	user := c.Value("current_user").(models.User)
	q := tx.Q()
	status := "active"
	q.Where("Active = ?", status)
	if err := q.All(&users); err != nil {
		return err
	}
	UserList := []map[string]interface{}{}
	for _, user := range users {
		User := map[string]interface{}{
			user.Name + " " + user.LastName: user.ID,
		}
		UserList = append(UserList, User)
	}
	c.Set("usersList", UserList)
	c.Set("users", users)
	c.Set("user", user)
	c.Set("tasks", models.Task{})
	return c.Render(http.StatusOK, r.HTML("tasks/new.plush.html"))
}

func Newtaskuser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	users := models.Users{}
	user := c.Value("current_user").(models.User)
	q := tx.Q()
	status := "true"
	q.Where("Active = ?", status)
	if err := q.All(&users); err != nil {
		return err
	}
	c.Set("user", user)
	c.Set("users", users)
	c.Set("tasks", models.Task{})
	return c.Render(http.StatusOK, r.HTML("tasks/newtask.plush.html"))
}
func Createtask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := models.Task{}

	if err := c.Bind(&tasks); err != nil {
		return err
	}
	users := models.Users{}
	q := tx.Q()
	status := "true"
	q.Where("Active = ?", status)
	if err := q.All(&users); err != nil {
		return err
	}
	UserList := []map[string]interface{}{}
	for _, user := range users {
		User := map[string]interface{}{
			user.Name + " " + user.LastName: user.ID,
		}
		UserList = append(UserList, User)
	}

	verrs := tasks.ValidateCreate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("tasks", tasks)
		c.Set("usersList", UserList)

		return c.Render(http.StatusUnprocessableEntity, r.HTML("tasks/new.plush.html"))
	}
	if err := tx.Create(&tasks); err != nil {
		return err
	}

	c.Flash().Add("success", "task created success")
	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func Createtaskuser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := models.Task{}

	if err := c.Bind(&tasks); err != nil {
		return err
	}
	users := models.Users{}
	user := c.Value("current_user").(models.User)
	tasks.UserID = user.ID
	status := "true"
	q := tx.Where("Id= ?", user.ID).Where("Active = ?", status)
	if err := q.All(&users); err != nil {
		return err
	}

	verrs := tasks.Validate()
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("tasks", tasks)

		return c.Render(http.StatusUnprocessableEntity, r.HTML("tasks/newtask.plush.html"))
	}
	if err := tx.Create(&tasks); err != nil {
		return err
	}

	c.Flash().Add("success", "task created success")
	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func Showtask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	task := models.Task{}
	user := c.Value("current_user").(models.User)
	taskid := c.Param("task_id")
	if err := tx.Find(&task, taskid); err != nil {
		return c.Render(http.StatusNotFound, r.HTML("/tasks"))
	}
	c.Set("user", user)
	c.Set("task", task)
	return c.Render(http.StatusOK, r.HTML("tasks/showtask.plush.html"))
}

func Edittask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := models.Task{}
	taskid := c.Param("task_id")
	user := c.Value("current_user").(models.User)

	if err := tx.Find(&tasks, taskid); err != nil {
		return err
	}

	users := models.Users{}

	if err := tx.All(&users); err != nil {
		return err
	}
	UserList := []map[string]interface{}{}
	for _, user := range users {
		User := map[string]interface{}{
			user.Name + " " + user.LastName: user.ID,
		}
		UserList = append(UserList, User)
	}
	c.Set("usersList", UserList)
	c.Set("users", users)
	c.Set("user", user)
	c.Set("task", tasks)
	return c.Render(http.StatusOK, r.HTML("tasks/edit.plush.html"))
}

func Updatetask(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	task := models.Task{}
	taskid := c.Param("task_id")

	if err := tx.Find(&task, taskid); err != nil {
		return err
	}

	if err := c.Bind(&task); err != nil {
		return err
	}
	verrs := task.ValidateUpdate()
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("task", task)

		return c.Render(http.StatusOK, r.HTML("tasks/edit.plush.html"))
	}
	if err := tx.Update(&task); err != nil {
		return err
	}
	c.Flash().Add("primary", "Task updated success")

	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func Edittaskuser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := models.Task{}
	user := c.Value("current_user").(models.User)

	taskid := c.Param("task_id")

	if err := tx.Find(&tasks, taskid); err != nil {
		return err
	}
	tasks.UserID = c.Value("current_user").(models.User).ID
	c.Set("user", user)
	c.Set("task", tasks)
	return c.Render(http.StatusOK, r.HTML("tasks/edituser.plush.html"))
}

func Updatetaskuser(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	task := models.Task{}
	taskid := c.Value("current_user").(models.User).ID

	if err := tx.Find(&task, taskid); err != nil {
		return err
	}

	if err := c.Bind(&task); err != nil {
		return err
	}
	verrs := task.ValidateUpdate()
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("task", task)

		return c.Render(http.StatusOK, r.HTML("tasks/edituser.plush.html"))
	}
	if err := tx.Update(&task); err != nil {
		return err
	}
	c.Flash().Add("primary", "Task updated success")

	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func Delete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	taskid := c.Param("task_id")
	taskid2, _ := uuid.FromString(taskid)
	if taskid == "" {
		return c.Redirect(http.StatusNotFound, "/tasks")
	}
	tasks := models.Tasks{}
	taskdadmin := c.Value("current_user").(models.User)
	if taskdadmin.Rol == "user" {
		if err := tx.Where("user_id = ? ", &taskdadmin.ID).Where("id = ?", &taskid).All(&tasks); err != nil {
			return err
		}
		if taskdadmin.ID == taskid2 {
			if err := tx.Destroy(&taskdadmin); err != nil {
				return err
			}
			c.Flash().Add("danger", "Task delete success")
			return c.Redirect(http.StatusSeeOther, "/tasks")
		}
	}

	task := &models.Task{ID: taskid2}
	if err := tx.Destroy(task); err != nil {
		return err
	}
	c.Flash().Add("danger", "Task delete success")

	return c.Redirect(http.StatusSeeOther, "/tasks")
}

func Updatecomplete(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	task := models.Task{}
	taskid := c.Param("task_id")

	if err := tx.Find(&task, taskid); err != nil {
		return err
	}

	if err := c.Bind(&task); err != nil {
		return err
	}
	var current string
	if !task.Complete {
		task.Complete = true
		task.Date = time.Now()
		c.Flash().Add("primary", "Task completed success")
		current = "/tasks"
	} else if task.Complete {
		task.Date = time.Now()
		task.Complete = false
		c.Flash().Add("primary", "Task returned success")

		current = "/tasks"
	}
	if err := tx.Update(&task); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, current)
}
