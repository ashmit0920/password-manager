package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	m := initialModel()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
