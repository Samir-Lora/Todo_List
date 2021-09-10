package actions_test

import (
	"strings"
	"time"
	"todo_list/app/models"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (as *ActionSuite) Test_UserList() {
	//create models
	user := models.Users{{Name: "prueba", LastName: "Prueba", Email: "hola@gmail.com", Active: "inactive", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "prueba", LastName: "Prueba2", Email: "hola2@gmail.com", Active: "active", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}}
	for _, t := range user {
		t.Email = strings.ToLower(t.Email)
		ph, _ := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
		t.PasswordHash = string(ph)
		err := as.DB.Create(&t)
		as.NoError(err)
		as.Session.Set("current_user_id", t.ID)

	}

	//testing url "/" is a index

	res := as.HTML("/users").Get()

	body := res.Body.String()

	for _, t := range user {
		as.Contains(body, t.Email)

	}
	//response stauts
	as.Equal(200, res.Code)

}

func (as *ActionSuite) Test_NewUser() {
	//testing url "/user/new" is a index
	restask := as.HTML("/user/new").Get()
	//response stauts
	as.Equal(200, restask.Code)
}

func (as *ActionSuite) Test_CreateUser() {
	//testing url "/user/create" is post in new
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)

	user := models.User{ID: id, Email: "hola@gmail.com", Password: "hola", Name: "prueba", LastName: "Prueba", Active: "inactive", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	user.Email = strings.ToLower(user.Email)
	ph, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.PasswordHash = string(ph)

	res2 := as.HTML("/user/create").Post(&user)
	as.Equal(303, res2.Code)
	as.Equal("/", res2.Location())
	err2 := as.DB.First(&user)
	as.NoError(err2)
	as.NotZero(user.ID)
	as.NotZero(user.CreatedAt)
	//verification
	as.Equal("prueba", user.Name)
	as.Equal(id, user.ID)
}

func (as *ActionSuite) Test_Edituser() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	user := models.User{ID: id, Email: "hola@gmail.com", Active: "active", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	verrs, _ := as.DB.ValidateAndCreate(&user)
	as.Session.Set("current_user_id", user.ID)
	as.False(verrs.HasAny())
	err := as.DB.Reload(&user)
	as.NoError(err)

	resuser := as.HTML("/user/edit/%s", user.ID).Get()
	//response stauts
	as.Equal(200, resuser.Code)
}

func (as *ActionSuite) Test_Updateuser() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	user := models.User{ID: id, Email: "hola@gmail.com", Password: "hola", Name: "prueba", LastName: "Prueba", Active: "inactive", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	user.Email = strings.ToLower(user.Email)
	ph, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.PasswordHash = string(ph)
	verrs, _ := as.DB.ValidateAndCreate(&user)
	as.False(verrs.HasAny())
	as.Session.Set("current_user_id", user.ID)
	userupdate := models.User{ID: id, Email: "hola@gmail.com", Password: "hola", Name: "prueba2", LastName: "Prueba", Active: "inactive", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	userupdate.Email = strings.ToLower(user.Email)
	ph, _ = bcrypt.GenerateFromPassword([]byte(userupdate.Password), bcrypt.DefaultCost)

	userupdate.PasswordHash = string(ph)

	res := as.HTML("/user/edit/%s/update/", user.ID).Put(&userupdate)
	as.Equal(303, res.Code)

	err := as.DB.Reload(&user)
	as.NoError(err)
	as.Equal("prueba2", user.Name)

}

func (as *ActionSuite) Test_Deleteuser() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	user := models.User{ID: id, Email: "hola@gmail.com", Active: "active", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	verrs, _ := as.DB.ValidateAndCreate(&user)
	as.Session.Set("current_user_id", user.ID)

	as.False(verrs.HasAny())

	res := as.HTML("/user/delete/{%s}", user.ID).Delete()
	as.Equal(303, res.Code)
}

func (as *ActionSuite) Test_Updateactive() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	user := models.User{ID: id, Email: "hola@gmail.com", Active: "active", Rol: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	verrs, _ := as.DB.ValidateAndCreate(&user)
	as.False(verrs.HasAny())
	as.Session.Set("current_user_id", user.ID)

	if user.Active == "active" {
		res := as.HTML("/user/updateactive/%s", user.ID).Put(&models.User{ID: id, Name: "prueba2", LastName: "Prueba", Email: "hola@gmail.com", Active: "inactive", CreatedAt: time.Now(), UpdatedAt: time.Now()})
		as.Equal(303, res.Code)
		err := as.DB.Reload(&user)
		as.NoError(err)
		as.Equal("prueba2", user.Name)

	} else if user.Active == "inactive" {
		res := as.HTML("/user/updateactive/%s", user.ID).Put(&models.User{ID: id, Name: "prueba3", LastName: "Prueba", Email: "hola@gmail.com", Active: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()})
		as.Equal(303, res.Code)
		err := as.DB.Reload(&user)
		as.NoError(err)
		as.Equal("prueba3", user.Name)
	}

}
