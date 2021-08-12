package actions

import (
	"net/http"
	"todo_list/app/models"
	"todo_list/app/render"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5" // modificado
	"github.com/gofrs/uuid"
)

var (
	// r is a buffalo/render Engine that will be used by actions
	// on this package to render render HTML or any other formats.
	r = render.Engine
)

func TaskList(i buffalo.Context) error {
	tx := i.Value("tx").(*pop.Connection)

	status := i.Param("complete")
	q := tx.Q()
	if status != "" {
		q.Where("Complete = ?", status)
	}

	tasks := models.Tasks{}

	if err := q.All(&tasks); err != nil {
		return err
	}

	i.Set("tasks", tasks)
	i.Set("len", len(tasks))

	return i.Render(http.StatusOK, r.HTML("home/incompleted.plush.html"))
}

func Newtask(i buffalo.Context) error {
	i.Set("tasks", models.Task{})
	return i.Render(http.StatusOK, r.HTML("home/new.plush.html"))
}

func Createtask(i buffalo.Context) error {
	tx := i.Value("tx").(*pop.Connection)
	tasks := models.Task{}

	if err := i.Bind(&tasks); err != nil {
		return err
	}

	verrs := tasks.Validate()
	if verrs.HasAny() {
		i.Set("errors", verrs)
		i.Set("tasks", tasks)
		i.Set("task", tasks.Task)
		i.Set("Description", tasks.Description)
		i.Set("Date", tasks.Date)

		return i.Render(http.StatusOK, r.HTML("home/new.plush.html"))
	}

	if err := tx.Create(&tasks); err != nil {
		return err
	}

	return i.Redirect(http.StatusSeeOther, "/task")
}

func Showtask(i buffalo.Context) error {
	tx := i.Value("tx").(*pop.Connection)
	task := models.Task{}
	taskID := i.Param("task_ID")

	if err := tx.Find(&task, taskID); err != nil {
		return err
	}

	i.Set("task", task)
	return i.Render(http.StatusOK, r.HTML("home/showtask.plush.html"))
}

func Edittask(i buffalo.Context) error {
	tx := i.Value("tx").(*pop.Connection)
	tasks := models.Task{}
	taskID := i.Param("task_ID")

	if err := tx.Find(&tasks, taskID); err != nil {
		return err
	}

	i.Set("task", tasks)
	return i.Render(http.StatusOK, r.HTML("home/edit.plush.html"))
}

func Updatetask(i buffalo.Context) error {

	tx := i.Value("tx").(*pop.Connection)
	task := models.Task{}
	taskID := i.Param("task_ID")

	if err := tx.Find(&task, taskID); err != nil {
		return err
	}

	if err := i.Bind(&task); err != nil {
		return err
	}
	verrs := task.Validate()
	if verrs.HasAny() {
		i.Set("errors", verrs)
		i.Set("task", task)

		return i.Render(http.StatusOK, r.HTML("home/edit.plush.html"))
	}
	if err := tx.Update(&task); err != nil {
		return err
	}
	return i.Redirect(http.StatusSeeOther, "/task?complete=false")
}

func Delete(i buffalo.Context) error {
	tx := i.Value("tx").(*pop.Connection)
	taskID := i.Param("task_ID")
	taskID2, _ := uuid.FromString(taskID)
	tasks := &models.Task{ID: taskID2}
	if err := tx.Destroy(tasks); err != nil {
		return err
	}

	return i.Redirect(http.StatusSeeOther, "/task")
}

func Updatecomplete(i buffalo.Context) error {

	tx := i.Value("tx").(*pop.Connection)
	task := models.Task{}
	taskID := i.Param("task_ID")

	if err := tx.Find(&task, taskID); err != nil {
		return err
	}

	if err := i.Bind(&task); err != nil {
		return err
	}
	var current string
	if !task.Complete {
		task.Complete = true
		current = "/task?complete=false"
	} else if task.Complete {
		task.Complete = false
		current = "/task?complete=true"
	}
	if err := tx.Update(&task); err != nil {
		return err
	}
	return i.Redirect(http.StatusSeeOther, current)
}
