package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	width, height int
	quiting       bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			m.quiting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var view tea.View
	if m.quiting {
		return view
	}
	return view
}

func main() {
	fmt.Println("Hello Mom")
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(69)
	}
}
