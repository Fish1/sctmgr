package gemgrscreen

func (m Model) View() string {
	s := m.list.View()
	return s
}
