package actions

import (
	"fmt"
	"net/http"
	"strings"
	"todo_list/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

//Render the list of users
func UserList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(models.User)
	users := models.Users{}
	q := tx.PaginateFromParams(c.Params())
	if err := q.Order("active, rol, name").All(&users); err != nil {
		return err
	}
	c.Set("userspaginator", q.Paginator)
	c.Set("user", user)
	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("users/userslist.plush.html"))
}

//NewUser renders the users form
func NewUser(c buffalo.Context) error {
	c.Set("users", models.User{})
	return c.Render(http.StatusOK, r.HTML("users/newuser.plush.html"))
}

//Validate and create the new user
func CreateUser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return errors.WithStack(err)
	}

	//create a password hash
	err := user.Create(tx)
	if err != nil {
		return err
	}

	//validate a values

	verrs := user.Validate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("users", user)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("users/newuser.plush.html"))
	}
	//set a necessary params

	user.Rol = "user"
	user.Active = "active"
	user.Email = strings.ToLower(user.Email)

	//Create user in database
	if err := tx.Create(&user); err != nil {
		return err
	}

	c.Flash().Add("success", "User created success")
	return c.Redirect(http.StatusSeeOther, "/")
}

//NewInvitation renders the users form
func NewInvitation(c buffalo.Context) error {
	c.Set("user", models.User{})

	return c.Render(http.StatusOK, r.HTML("users/newinvitation.plush.html"))
}

// create a new user with role invited by admin
func CreateInvitation(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		return err
	}
	verrs := user.ValidateCreate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("users", user)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("users/newinvitation.plush.html"))
	}
	user.Active = "invited"
	user.Email = strings.ToLower(user.Email)
	if err := tx.Create(&user); err != nil {
		return err
	}

	c.Flash().Add("success", "User created success")
	return c.Redirect(http.StatusSeeOther, "/users")
}

// render html with form for the new password as the guest user
func CreatePassByInvitation(c buffalo.Context) error {
	c.Set("user", &models.User{})
	return c.Render(http.StatusOK, r.HTML("users/newpassbyinvitation.plush.html"))
}

//Create a passoword the guest user
func UpdatePassByInvitation(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userid := c.Param("user_id")
	if err := tx.Find(&user, userid); err != nil {
		return err
	}

	if err := c.Bind(&user); err != nil {
		return err
	}
	//create a password hash
	err := user.Create(tx)
	if err != nil {
		return err
	}
	user.Active = "active"

	verrs := user.Validate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)

		return c.Render(http.StatusOK, r.HTML("users/newpassbyinvitation.plush.html"))
	}
	if err := tx.Update(&user); err != nil {
		return err
	}
	c.Flash().Add("primary", "User updated success")

	return c.Redirect(http.StatusSeeOther, "/")
}

//Render a html with form for update password
func Updatepassword(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userid := c.Value("current_user").(models.User).ID

	if err := tx.Find(&user, userid); err != nil {
		return err
	}
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/changepassword.plush.html"))
}

//Function change password of current user
func Changepass(c buffalo.Context) error {
	userid := c.Param("user_id")
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(models.User)

	if err := tx.Find(&user, userid); err != nil {
		return err
	}

	if err := c.Bind(&user); err != nil {
		return err
	}
	//create a password hash
	err := user.Create(tx)
	if err != nil {
		return err
	}
	verrs := user.ValidateUpdatePassword(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)

		return c.Render(http.StatusOK, r.HTML("users/changepassword.plush.html"))
	}
	if err := tx.Update(&user); err != nil {
		return err
	}
	c.Session().Clear()
	c.Flash().Add("primary", "User password updated success")
	return c.Redirect(http.StatusSeeOther, "/")
}

//render list of user
func Showuser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(models.User)
	userid := c.Param("user_id")
	if user.Rol == "user" && user.ID.String() == userid {
		if err := tx.Find(&user, userid); err != nil {
			fmt.Println("holi")

			c.Set("user", user)
			return c.Render(http.StatusOK, r.HTML("users/showuser.plush.html"))
		}
	} else if user.Rol == "admin" {
		if err := tx.Find(&user, userid); err != nil {
			return c.Render(http.StatusNotFound, r.HTML("/users"))
		}
	} else {
		c.Flash().Add("danger", "You must be authorized to see that page")
		return c.Redirect(302, "/tasks")
	}

	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("users/showuser.plush.html"))
}

//Render html for edit user
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

//update a user by admin
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
	verrs := user.ValidateUpdate()
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

//delete user by admin
func Deleteuser(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	userid := c.Param("user_id")
	userid2, _ := uuid.FromString(userid)
	if userid == "" {
		return c.Redirect(http.StatusNotFound, "/users")
	}

	//delete current user
	rol := "admin"
	user := c.Value("current_user").(models.User)
	usersadmin := &models.Users{}
	if err := tx.Where("id = ?", user.ID).Where("rol =?", rol).All(usersadmin); err != nil {
		return err
	}

	if user.Rol == "admin" && user.ID == userid2 {
		if err := tx.Destroy(usersadmin); err != nil {
			return err
		}
		c.Session().Clear()
		c.Flash().Add("danger", "User delete success")
		return c.Redirect(http.StatusSeeOther, "/")
	}

	//delete user or admin
	users := &models.User{ID: userid2}
	if err := tx.Destroy(users); err != nil {
		return err
	}
	c.Flash().Add("danger", "User delete success")

	return c.Redirect(http.StatusSeeOther, "/users")
}

//update active or desactive status user by admin
func Updateactive(c buffalo.Context) error {

	tx := c.Value("tx").(*pop.Connection)
	user := models.User{}
	userid := c.Param("user_id")
	currentuser := c.Value("current_user").(models.User)
	if err := tx.Find(&user, userid); err != nil {
		return err
	}

	if err := c.Bind(&user); err != nil {
		return err
	}

	var current string
	if user.Active == "active" {
		user.Active = "inactive"
		if currentuser.ID.String() == userid {
			c.Session().Clear()
			c.Flash().Add("success", "You have been logged out and inactive !")
			if err := tx.Update(&user); err != nil {
				return err
			}
			return c.Redirect(302, "/")
		}
		c.Flash().Add("primary", "User actived success")
		current = "/users"
	} else if user.Active == "inactive" {
		user.Active = "active"
		c.Flash().Add("primary", "User desactive success")

		current = "/users"
	}
	if err := tx.Update(&user); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, current)
}
