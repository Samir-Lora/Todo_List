package render_test

import (
	"testing"
	"time"
	"todo_list/app/render"
)

func Test_Date(t *testing.T) {
	Date := time.Date(2018, 01, 02, 0, 0, 0, 0, time.UTC)
	timeformat := render.Date(Date)
	if timeformat == "" {
		t.Error("This date is empty")
	}
	if timeformat != "02 Jan 2018" {
		t.Error("The date excepted", timeformat)
	}
}

func Test_Datecompleted(t *testing.T) {
	status := true
	actual := render.Datecompleted(status)
	if actual != "Was completed on" {
		t.Error("Unconverted, value is : ", actual)
	}
	status2 := false
	actual2 := render.Datecompleted(status2)
	if actual2 != "Needs to be completed on" {
		t.Error("Unconverted, value is : ", actual2)
	}

}

func Test_SelectedFilter(t *testing.T) {
	location := render.SelectedFilter("true", "true")
	if location != "font-weight-lighter" {
		t.Error("Unconverted, value is : ", location)
	}
}

func Test_Button(t *testing.T) {
	location := render.SelectedFilter("true", "true")
	if location != "d" {
		t.Error("Unconverted, value is : ", location)
	}
}
