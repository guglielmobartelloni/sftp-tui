package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := model{
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
	err := tea.NewProgram(m, tea.WithAltScreen()).Start()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	connectToServer()

}

type model struct {
	stopwatch stopwatch.Model
}

func (m model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd

}

func (m model) View() string {
	s := m.stopwatch.View() + "\n"
	s = "Elapsed: " + s
	return s
}
