package menu

import (
	"github.com/charmbracelet/bubbles/list"
)

type Choice struct {
	name        string
	downloadUrl string
}

type Selection struct {
	status   Status
	localUri string
}

type Model struct {
	list list.Model
}
