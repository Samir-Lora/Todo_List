package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Task is used by pop to map your Tasks database table to your go code.
type Task struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Task        string    `json:"task" db:"task"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	Complete    bool      `json:"complete" db:"complete"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (i Tasks) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Tasks is not required by pop and may be deleted
type Tasks []Task

// String is not required by pop and may be deleted

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.

func (i *Task) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Task, Name: "Task"},
		&validators.StringIsPresent{Field: i.Description, Name: "Description"},
		&validators.TimeIsPresent{Field: i.Date, Name: "Date"},
	)
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *Task) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *Task) ValidateUpdate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: i.Task, Name: "Task"},
		&validators.StringIsPresent{Field: i.Description, Name: "Description"},
		&validators.TimeIsPresent{Field: i.Date, Name: "Date"},
	)
}
