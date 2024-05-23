package menu

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Status int64

const (
	None        Status = 0
	Downloading Status = 1
	Downloaded  Status = 2
	Deleteing   Status = 3
)

type GEItem struct {
	name       string
	localPath  string
	remotePath string
	status     Status
}

type item GEItem

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	cursor := " "
	if index == m.Index() {
		cursor = ">"
	}

	var status string
	switch i.status {
	case None:
		status = " "
	case Downloaded:
		status = "X"
	case Downloading:
		status = "D"
	case Deleteing:
		status = "d"
	}

	str := fmt.Sprintf("%s [%s] %s\t\t\t%s", cursor, status, i.name, i.localPath)
	fmt.Fprintf(w, "%s", str)
}
