// middleware package is intended to host the middlewares used
// across the app.
package middleware

import (
	"net/http"
	"time"
	"todo_list/app/models"

	"github.com/gobuffalo/buffalo"
	tx "github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
)

var (
	// Transaction middleware wraps the request with a pop
	// transaction that is committed on success and rolled
	// back when errors happen.
	Transaction = tx.Transaction(models.DB())

	// ParameterLogger logs out parameters that the app received
	// taking care of sensitive data.
	ParameterLogger = paramlogger.ParameterLogger

	// CSRF middleware protects from CSRF attacks.
	CSRF = csrf.New
)

func IncompleteTask(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tx := models.DB()
		q := tx.Q()
		q.Where("Complete = false")
		tasks := models.Tasks{}
		if err := q.All(&tasks); err != nil {
			return err
		}
		c.Set("len", len(tasks))
		return next(c)
	}
}

func EditTaskAcess(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tx := models.DB()
		task := models.Task{}
		taskid := c.Param("task_id")
		tx.Find(&task, taskid)
		if task.Complete {
			c.Flash().Add("danger", "cannot edit a complete task")
			c.Redirect(http.StatusSeeOther, "/")
		}
		return next(c)
	}

}

func Datenow(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		currentTime := time.Now().Format("Monday 02, January 2006")

		c.Set("datenow", currentTime)

		return next(c)
	}
}
