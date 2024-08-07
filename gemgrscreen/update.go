package gemgrscreen

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fish1/sctmgr/gemgr"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd := handleKeyMsg(m, msg)
		if cmd != nil {
			return m, cmd
		}
	case spinner.TickMsg:
		return handleSpinnerTick(m, msg)
	case DoneDelete:
		return handleDoneDelete(m, msg)
	case DoneDownload:
		return handleDoneDownload(m, msg)
	case FailedDownload:
		return handleFailedDownload(m, msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func handleSpinnerTick(m Model, msg spinner.TickMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	m.list.SetDelegate(itemDelegate{
		spinner: &m.spinner,
	})
	return m, cmd
}

func handleKeyMsg(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case " ", "enter":
		i, ok := m.list.SelectedItem().(item)
		if !ok {
			panic(ok)
		}
		if i.status == Downloaded {
			return handleStartDelete(m, i)
		} else if i.status == None {
			m, cmd := handleStartDownload(m, i)
			return m, cmd
		} else if i.status == Downloading {
			i.cancel()
		}
	}
	return m, nil
}

func handleStartDelete(m Model, i item) (tea.Model, tea.Cmd) {
	i.status = Deleteing
	i.cancel = nil
	m.list.SetItem(m.list.Index(), item(i))
	return m, func() tea.Msg {
		err := gemgr.Delete(i.localPath)
		if err != nil {
			panic(err)
		}
		return DoneDelete{index: m.list.Index()}
	}
}

func handleDoneDelete(m Model, msg DoneDelete) (tea.Model, tea.Cmd) {
	i, ok := m.list.Items()[msg.index].(item)
	if !ok {
		panic(ok)
	}
	i.status = None
	i.localPath = ""
	i.cancel = nil
	cmd := m.list.SetItem(msg.index, item(i))
	return m, cmd
}

func handleStartDownload(m Model, i item) (tea.Model, tea.Cmd) {
	ctx, cancel := context.WithCancel(context.Background())
	i.status = Downloading
	i.cancel = cancel
	m.list.SetItem(m.list.Index(), item(i))
	return m, func() tea.Msg {
		localPath, err := gemgr.Install(i.remotePath, ctx)
		if err != nil {
			return FailedDownload{index: m.list.Index()}
		}
		return DoneDownload{index: m.list.Index(), localPath: localPath}
	}
}

func handleDoneDownload(m Model, msg DoneDownload) (tea.Model, tea.Cmd) {
	i, ok := m.list.Items()[msg.index].(item)
	if !ok {
		panic(ok)
	}
	i.status = Downloaded
	i.localPath = msg.localPath
	i.cancel = nil
	cmd := m.list.SetItem(msg.index, item(i))
	return m, cmd
}

func handleFailedDownload(m Model, msg FailedDownload) (tea.Model, tea.Cmd) {
	i, ok := m.list.Items()[msg.index].(item)
	if !ok {
		panic(ok)
	}
	i.status = None
	i.localPath = ""
	i.cancel = nil
	cmd := m.list.SetItem(msg.index, item(i))
	return m, cmd
}

type DoneDownload struct {
	index     int
	localPath string
}

type DoneDelete struct {
	index int
}

type FailedDownload struct {
	index int
}
