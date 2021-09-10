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
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
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
		tasks := models.Tasks{}
		uid := c.Session().Get("current_user_id")
		if uid != nil {
			user := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(user, uid)
			if err != nil {
				return err
			}
			if user.Rol == "admin" {
				q.Where("complete = false")
				if err := q.All(&tasks); err != nil {
					return err
				}
			} else {
				q.Where("complete = false").Where("user_id = ?", uid)
				if err := q.All(&tasks); err != nil {
					return err
				}
			}
			c.Set("rol", user.Rol)
			c.Set("len", len(tasks))
		}
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
			c.Flash().Add("danger", "Cannot edit a complete task")
			c.Redirect(http.StatusSeeOther, "/tasks")
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

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(&u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}

func Authorizeusers(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user := c.Value("current_user").(models.User)
		if user.Rol == "admin" {
			return next(c)
		}
		c.Flash().Add("danger", "You must be authorized to see that page")
		return c.Redirect(302, "/tasks")
	}
}
