package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fish1/sctmgr/gemgrscreen"
	// "github.com/fish1/sctmgr/mainscreen"
)

func main() {
	program := tea.NewProgram(gemgrscreen.New())
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}
