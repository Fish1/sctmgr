package menu

import "github.com/charmbracelet/bubbles/spinner"

type Status int64

const (
	Idle     Status = 0
	Download Status = 1
	Delete   Status = 2
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
	header   string
	choices  []Choice
	cursor   int
	selected map[int]Selection
	spinner  spinner.Model
}
