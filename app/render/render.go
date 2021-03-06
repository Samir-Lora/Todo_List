package render

import (
	"time"
	base "todo_list"

	"github.com/gobuffalo/buffalo/render"
)

// Engine for rendering across the app, it provides
// the base for rendering HTML, JSON, XML and other formats
// while also defining thing like the base layout.
var Engine = render.New(render.Options{
	HTMLLayout:   "application.plush.html",
	TemplatesBox: base.Templates,
	AssetsBox:    base.Assets,
	Helpers:      Helpers,
})

// Helpers available for the plush templates, there are
// some helpers that are injected by Buffalo but this is
// the list of custom Helpers.
var Helpers = map[string]interface{}{
	// partialFeeder is the helper used by the render engine
	// to find the partials that will be used, this is important
	"partialFeeder":  base.Templates.FindString,
	"Date":           Date,
	"Datecompleted":  Datecompleted,
	"SelectedFilter": SelectedFilter,
	"Button":         Button,
	"Priority":       Priority,
	"Prioritycolor":  Prioritycolor,
	"Priorityicon":   Priorityicon,
}

func Date(value time.Time) string {
	date := value.Format("02 Jan 2006")
	return date
}

func Datecompleted(statusref bool) string {
	var status string
	if statusref {
		status = "Was completed on"
	} else {
		status = "Needs to be completed on"
	}
	return status
}

func SelectedFilter(completed string, expectedValue string) string {
	var value string
	if completed == expectedValue {
		value = "font-weight-lighter"
	}
	return value

}

func Button(completed string, expectedValue string) string {
	var value string
	if completed == expectedValue {
		value = "d"
	}
	return value
}

func Priority(status string) string {

	var statuspriority string
	if status == "1" {
		statuspriority = "High"
	} else if status == "2" {
		statuspriority = "Medium"
	} else if status == "3" {
		statuspriority = "Low"
	}
	return statuspriority
}

func Prioritycolor(status string) string {

	var statuspriority string

	if status == "1" {
		statuspriority = "text-danger"
	} else if status == "2" {
		statuspriority = "text-warning"
	} else if status == "3" {
		statuspriority = "text-success"
	}
	return statuspriority
}

func Priorityicon(status string) string {
	var iconpriority string
	if status == "1" {
		iconpriority = "exclamation-circle"
	} else if status == "2" {
		iconpriority = "exclamation-triangle"
	} else if status == "3" {
		iconpriority = "info-circle"
	}
	return iconpriority
}
