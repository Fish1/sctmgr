package menu

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
)

type Model struct {
	list    list.Model
	spinner spinner.Model
}
