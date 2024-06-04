package controllers

import (
	"raiden/internal/models"

	"github.com/sev-2/raiden"
)

type ToDoController struct {
	raiden.ControllerBase
	Http  string `path:"/todo-list" type:"rest"`
	Model models.List
}
