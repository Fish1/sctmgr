package menu

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fish1/ge-downloader/gemanager"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKeyMsg(m, msg)
	case spinner.TickMsg:
		return handleSpinnerTick(m, msg)
	case DoneDelete:
		return handleDoneDelete(m, msg)
	case DoneDownload:
		return handleDoneDownload(m, msg)
	default:
		return m, nil
	}
}

func handleKeyMsg(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case " ", "enter":
		if s, ok := m.selected[m.cursor]; ok {
			s.status = Delete
			return m, func() tea.Msg {
				err := gemanager.Delete(s.localUri)
				if err != nil {
					panic(err)
				}
				return DoneDelete{index: m.cursor}
			}
		} else {
			m.selected[m.cursor] = Selection{
				status: Download,
			}
			return m, func() tea.Msg {
				filename, err := gemanager.Install(m.choices[m.cursor].downloadUrl)
				panic(filename)
				if err != nil {
					panic(err)
				}
				return DoneDownload{index: m.cursor, filename: filename}
			}
		}
	}
	return m, nil
}

func handleSpinnerTick(m Model, msg spinner.TickMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func handleDoneDelete(m Model, msg DoneDelete) (tea.Model, tea.Cmd) {
	delete(m.selected, msg.index)
	return m, nil
}

func handleDoneDownload(m Model, msg DoneDownload) (tea.Model, tea.Cmd) {
	m.selected[msg.index] = Selection{
		status:   Idle,
		localUri: msg.filename,
	}
	return m, nil
}

type DoneDownload struct {
	index    int
	filename string
}

type DoneDelete struct {
	index int
}
