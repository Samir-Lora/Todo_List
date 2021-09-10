package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is used by pop to map your users database table to your go code.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	LastName     string    `json:"lastname" db:"lastname"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`

	Password             string `json:"-" db:"-"`
	PasswordConfirmation string `json:"-" db:"-"`

	Active string `json:"active" db:"active"`
	Rol    string `json:"rol" db:"rol"`

	Tasks     Tasks     `has_many:"tasks"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) Create(tx *pop.Connection) error {
	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(ph)
	return err
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.

func (c *User) Validate(tx *pop.Connection) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.StringIsPresent{Field: c.LastName, Name: "LastName"},
		&validators.EmailIsPresent{Field: c.Email, Name: "Email"},
		&validators.StringIsPresent{Field: c.Password, Name: "Password"},
		&validators.StringIsPresent{Field: c.Password, Name: "PasswordConfirmation"},
		&validators.FuncValidator{
			Field:   c.Email,
			Name:    "Email",
			Message: "%s is already exist",
			Fn: func() bool {
				var exist bool
				q := tx.Where("email = ?", c.Email)
				if c.ID != uuid.Nil {
					q = q.Where("id != ?", c.ID)
				}
				exist, _ = q.Exists(c)
				return !exist
			},
		},
	)
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *User) ValidateCreate(tx *pop.Connection) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.StringIsPresent{Field: c.LastName, Name: "LastName"},
		&validators.EmailIsPresent{Field: c.Email, Name: "Email"},
		&validators.FuncValidator{
			Field:   c.Email,
			Name:    "Email",
			Message: "%s is already exist",
			Fn: func() bool {
				var exist bool
				q := tx.Where("email = ?", c.Email)
				if c.ID != uuid.Nil {
					q = q.Where("id != ?", c.ID)
				}
				exist, _ = q.Exists(c)
				return !exist
			},
		},
	)
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *User) ValidateUpdate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.StringIsPresent{Field: c.LastName, Name: "LastName"},
		&validators.EmailIsPresent{Field: c.Email, Name: "Email"},
	)
}

func (c *User) ValidateUpdatePassword(tx *pop.Connection) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Password, Name: "Password"},
		&validators.StringIsPresent{Field: c.Password, Name: "PasswordConfirmation"},
		&validators.FuncValidator{
			Name:    "Password",
			Message: "Can't have the same password, try a new %v",
			Fn: func() bool {
				var exist bool
				q := tx.Where("password_hash = ?", c.PasswordHash)
				if c.ID != uuid.Nil {
					q = q.Where("id != ?", c.ID)
				}
				exist, _ = q.Exists(c)
				return exist
			},
		},
	)
}
