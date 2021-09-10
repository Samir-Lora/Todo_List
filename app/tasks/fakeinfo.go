package tasks

import (
	"math/rand"
	"strings"
	"time"
	"todo_list/app/models"

	"github.com/gofrs/uuid"
	"github.com/markbates/grift/grift"
	"github.com/wawandco/fako"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Name         string `fako:"first_name"`
	LastName     string `fako:"last_name"`
	Email        string `fako:"email_address"`
	PasswordHash string `fako:"simple_password"`

	Password             string
	PasswordConfirmation string
}

type Task struct {
	Task        string `fako:"job_title"`
	Description string
	Date        time.Time
	Complete    bool
	Priority    string
	UserID      uuid.UUID
}

var _ = grift.Add("fakeinfouser", func(c *grift.Context) error {
	userfake := User{}
	tx := models.DB()
	for i := 0; i < 5; i++ {
		fako.Fill(&userfake)
		userfake.Email = strings.ToLower(userfake.Email)
		userfake.Password = "prueba"
		ph, _ := bcrypt.GenerateFromPassword([]byte(userfake.Password), bcrypt.DefaultCost)
		userfake.PasswordHash = string(ph)
		user := models.User{
			Email:        userfake.Email,
			Name:         userfake.Name,
			LastName:     userfake.LastName,
			PasswordHash: userfake.PasswordHash,
			Active:       "active",
			Rol:          "admin",
		}
		if err := tx.Create(&user); err != nil {
			return err
		}
	}

	for i := 0; i < 200; i++ {
		fako.Fill(&userfake)
		userfake.Email = strings.ToLower(userfake.Email)
		userfake.Password = "prueba"
		ph, _ := bcrypt.GenerateFromPassword([]byte(userfake.Password), bcrypt.DefaultCost)
		userfake.PasswordHash = string(ph)
		user := models.User{
			Email:        userfake.Email,
			Name:         userfake.Name,
			LastName:     userfake.LastName,
			PasswordHash: userfake.PasswordHash,
			Active:       "active",
			Rol:          "user",
		}
		if err := tx.Create(&user); err != nil {
			return err
		}
		taskfake := Task{}
		for i := 0; i < 5; i++ {
			fako.Fill(&taskfake)
			r := rand.Intn(2)
			if r == 0 {
				taskfake.Complete = true
			} else if r == 1 {
				taskfake.Complete = false
			}
			p := rand.Intn(3)
			if p == 0 {
				taskfake.Priority = "1"
			}
			if p == 1 {
				taskfake.Priority = "2"
			}
			if p == 2 {
				taskfake.Priority = "3"
			}
			task := models.Task{
				Task:        taskfake.Task,
				Description: "Test",
				Date:        time.Now(),
				Complete:    taskfake.Complete,
				Priority:    taskfake.Priority,
				UserID:      user.ID,
			}
			if err := tx.Create(&task); err != nil {
				return err
			}
		}
	}

	return nil
})
