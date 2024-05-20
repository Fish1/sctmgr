package menu

import "fmt"

func (m Model) View() string {
	s := fmt.Sprintf("%s\n\n", m.header)

	for i, c := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		path := " "
		status := " "
		if s, ok := m.selected[i]; ok {
			if s.status == Download || s.status == Delete {
				status = m.spinner.View()
			} else {
				status = "X"
			}
			path = s.localUri
		} else {
			path = c.downloadUrl
		}

		s += fmt.Sprintf("%s [%s] %s\t\t%s\n", cursor, status, c.name, path)
	}

	s += "\nPress q to quit.\n"
	return s
}
