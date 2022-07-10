package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	tea.NewProgram(&Model{}).Start()
}

type Model struct{}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return ""
}
