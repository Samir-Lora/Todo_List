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
	q := tx.Q()
	if status != "" {
		q.Where("Complete = ?", status)
	}
	tasks := models.Tasks{}
	if err := q.Order("date asc").All(&tasks); err != nil {
		return err
	}

	c.Set("tasks", tasks)

	return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
}

func Newtask(c buffalo.Context) error {
	c.Set("tasks", models.Task{})
	return c.Render(http.StatusOK, r.HTML("tasks/new.plush.html"))
}

func Createtask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := models.Task{}

	if err := c.Bind(&tasks); err != nil {
		return err
	}

	verrs := tasks.ValidateCreate()
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("tasks", tasks)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("tasks/new.plush.html"))
	}

	if err := tx.Create(&tasks); err != nil {
		return err
	}
	c.Flash().Add("success", "task created success")
	return c.Redirect(http.StatusSeeOther, "/")
}

func Showtask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	task := models.Task{}
	taskid := c.Param("task_id")
	if err := tx.Find(&task, taskid); err != nil {
		return c.Render(http.StatusNotFound, r.HTML("/"))
	}

	c.Set("task", task)
	return c.Render(http.StatusOK, r.HTML("tasks/showtask.plush.html"))
}

func Edittask(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := models.Task{}
	taskid := c.Param("task_id")

	if err := tx.Find(&tasks, taskid); err != nil {
		return err
	}

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

	return c.Redirect(http.StatusSeeOther, "/")
}

func Delete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	taskid := c.Param("task_id")
	taskid2, _ := uuid.FromString(taskid)
	if taskid == "" {
		return c.Redirect(http.StatusNotFound, "/")
	}
	tasks := &models.Task{ID: taskid2}
	if err := tx.Destroy(tasks); err != nil {
		return err
	}
	c.Flash().Add("danger", "Task delete success")

	return c.Redirect(http.StatusSeeOther, "/")
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
		current = "/"
	} else if task.Complete {
		task.Date = time.Now()
		task.Complete = false
		c.Flash().Add("primary", "Task returned success")

		current = "/"
	}
	if err := tx.Update(&task); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, current)
}
