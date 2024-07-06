package gemgrscreen

import (
	"context"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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
	cancel     context.CancelFunc
}

type item GEItem

func (i item) FilterValue() string { return i.name }

type itemDelegate struct {
	spinner *spinner.Model
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(item)
	if !ok {
		return
	}

	cursor := " "
	if index == m.Index() {
		if item.status == Downloaded || item.status == Downloading {
			cursor = "x"
		} else {
			cursor = ">"
		}
	}

	var status string
	switch item.status {
	case None:
		status = " "
	case Downloaded:
		status = "X"
	case Downloading:
		status = d.spinner.View()
	case Deleteing:
		status = d.spinner.View()
	}

	str := fmt.Sprintf("%s [%s] %s \t\t\t%s", cursor, status, item.name, item.localPath)
	fmt.Fprintf(w, "%s", str)
}
