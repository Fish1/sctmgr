package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fish1/sctmgr/menu"
)

func main() {
	program := tea.NewProgram(menu.New())
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}
