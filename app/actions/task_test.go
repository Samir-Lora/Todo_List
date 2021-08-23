package actions_test

import (
	"time"
	"todo_list/app/models"

	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_TaskList() {
	//create models
	task := models.Tasks{{Task: "prueba", Description: "Prueba", Date: time.Now(), Complete: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Task: "prueba", Description: "Prueba2", Date: time.Now(), Complete: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}}
	for _, t := range task {
		err := as.DB.Create(&t)
		as.NoError(err)
	}
	//testing url "/" is a index

	res := as.HTML("/").Get()
	res2 := as.HTML("/?complete=true").Get()
	res3 := as.HTML("/?complete=false").Get()

	body := res.Body.String()
	body2 := res2.Body.String()
	body3 := res3.Body.String()

	for _, t := range task {
		as.Contains(body, t.Task)
		as.Contains(body2, t.Task)
		as.Contains(body3, t.Task)
	}
	//response stauts
	as.Equal(200, res.Code)
	as.Equal(200, res2.Code)
	as.Equal(200, res3.Code)

}

func (as *ActionSuite) Test_Newtask() {
	//testing url "/task/new" is a index
	restask := as.HTML("/task/new").Get()
	//response stauts
	as.Equal(200, restask.Code)
}

func (as *ActionSuite) Test_Createtask() {
	//testing url "/task/create" is post in new
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	task := &models.Task{ID: id, Task: "prueba", Description: "Prueba", Date: time.Now(), Complete: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	res2 := as.HTML("/task/create").Post(task)
	as.Equal(303, res2.Code)
	as.Equal("/", res2.Location())
	err := as.DB.First(task)
	as.NoError(err)
	as.NotZero(task.ID)
	as.NotZero(task.CreatedAt)
	//verification
	as.Equal("prueba", task.Task)
	as.Equal(id, task.ID)
}

func (as *ActionSuite) Test_Edittask() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	task := &models.Task{ID: id, Task: "prueba", Description: "Prueba", Date: time.Now(), Complete: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	verrs, _ := as.DB.ValidateAndCreate(task)
	as.False(verrs.HasAny())
	err := as.DB.Reload(task)
	as.NoError(err)

	restask := as.HTML("/task/edit/%s", task.ID).Get()
	//response stauts
	as.Equal(200, restask.Code)
}

func (as *ActionSuite) Test_Updatetask() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	task := &models.Task{ID: id, Task: "prueba", Description: "Prueba", Date: time.Now(), Complete: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	verrs, _ := as.DB.ValidateAndCreate(task)
	as.False(verrs.HasAny())

	res := as.HTML("/task/edit/%s/update/", task.ID).Put(&models.Task{ID: task.ID, Task: "prueba2", Description: "Prueba2", Date: time.Now(), Complete: true, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	as.Equal(303, res.Code)

	err := as.DB.Reload(task)
	as.NoError(err)
	as.Equal("prueba2", task.Task)
	as.Equal("prueba2", task.Task)

}

func (as *ActionSuite) Test_Delete() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	task := &models.Task{ID: id, Task: "prueba", Description: "Prueba", Date: time.Now(), Complete: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	verrs, _ := as.DB.ValidateAndCreate(task)
	as.False(verrs.HasAny())

	res := as.HTML("/task/delete/{%s}", task.ID).Delete()
	as.Equal(303, res.Code)
}

func (as *ActionSuite) Test_Updatecomplete() {
	var x error
	id := uuid.Must(uuid.FromStringOrNil("2baaec43-8520-4120-8adf-c1f604fe30eb"), x)
	task := &models.Task{ID: id, Task: "prueba", Description: "Prueba", Date: time.Now(), Complete: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	verrs, _ := as.DB.ValidateAndCreate(task)
	as.False(verrs.HasAny())

	if task.Complete {
		res := as.HTML("/task/updatecomplete/%s", task.ID).Put(&models.Task{ID: task.ID, Task: "prueba2", Description: "Prueba2", Date: time.Now(), Complete: false, CreatedAt: time.Now(), UpdatedAt: time.Now()})
		as.Equal(303, res.Code)
		err := as.DB.Reload(task)
		as.NoError(err)
		as.Equal("prueba2", task.Task)

	} else if !task.Complete {
		res := as.HTML("/task/updatecomplete/%s", id).Put(&models.Task{ID: id, Task: "prueba3", Description: "Prueba2", Date: time.Now(), Complete: true, CreatedAt: time.Now(), UpdatedAt: time.Now()})
		as.Equal(303, res.Code)
		err := as.DB.Reload(task)
		as.NoError(err)
		as.Equal("prueba3", task.Task)
	}

}
