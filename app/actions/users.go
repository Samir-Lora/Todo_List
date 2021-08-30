package actions

import (
	"net/http"
	"todo_list/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

func UserList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	users := models.Users{}
	if err := tx.All(&users); err != nil {
		return err
	}
	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("users/userslist.plush.html"))
}

func NewUser(c buffalo.Context) error {
	c.Set("users", models.User{})
	return c.Render(http.StatusOK, r.HTML("users/newuser.plush.html"))
}

func CreateUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		return err
	}
	verrs := user.Validate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("users", user)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("users/newuser.plush.html"))
	}
	if err := tx.Create(&user); err != nil {
		return err
	}
	c.Flash().Add("success", "User created success")
	return c.Redirect(http.StatusSeeOther, "/users")
}

func Showuser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userid := c.Param("user_id")
	if err := tx.Find(&user, userid); err != nil {
		return c.Render(http.StatusNotFound, r.HTML("/tasks"))
	}

	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/showuser.plush.html"))
}

func Edituser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	users := models.User{}
	userid := c.Param("user_id")

	if err := tx.Find(&users, userid); err != nil {
		return err
	}

	c.Set("user", users)
	return c.Render(http.StatusOK, r.HTML("users/edit.plush.html"))
}

func Updateuser(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userid := c.Param("user_id")

	if err := tx.Find(&user, userid); err != nil {
		return err
	}

	if err := c.Bind(&user); err != nil {
		return err
	}
	verrs := user.Validate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)

		return c.Render(http.StatusOK, r.HTML("users/edit.plush.html"))
	}
	if err := tx.Update(&user); err != nil {
		return err
	}
	c.Flash().Add("primary", "User updated success")

	return c.Redirect(http.StatusSeeOther, "/users")
}

func Deleteuser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	userid := c.Param("user_id")
	userid2, _ := uuid.FromString(userid)
	if userid == "" {
		return c.Redirect(http.StatusNotFound, "/tasks")
	}
	users := &models.User{ID: userid2}
	if err := tx.Destroy(users); err != nil {
		return err
	}
	c.Flash().Add("danger", "User delete success")

	return c.Redirect(http.StatusSeeOther, "/users")
}

func Updatecompleteuser(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userid := c.Param("user_id")

	if err := tx.Find(&user, userid); err != nil {
		return err
	}

	if err := c.Bind(&user); err != nil {
		return err
	}
	var current string
	if err := tx.Update(&user); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, current)
}

func Updateactive(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userid := c.Param("user_id")

	if err := tx.Find(&user, userid); err != nil {
		return err
	}

	if err := c.Bind(&user); err != nil {
		return err
	}
	var current string
	if !user.Active {
		user.Active = true
		c.Flash().Add("primary", "User actived success")
		current = "/users"
	} else if user.Active {
		user.Active = false
		c.Flash().Add("primary", "User desactive success")

		current = "/users"
	}
	if err := tx.Update(&user); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, current)
}
