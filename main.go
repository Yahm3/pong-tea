

package main

import (
	"fmt"
	"os"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

type model struct {
	width, height    int
	quiting          bool
	ballX, ballY     int
	paddleX, paddleY int
	scoreA, scoreB   int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quiting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func playGround(str string) string {
	return lipgloss.NewStyle().
		Width(50).  // :TODO: make this dynamic
		Height(20). // :TODO: make this dynamic
		// Padding(1).
		PaddingBottom(0).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Charple).
		Align(lipgloss.Center, lipgloss.Center).
		Render(str)

}

func (m model) View() tea.View {
	var view tea.View

	view.AltScreen = true // immediately go into fullscreen mode

	if m.quiting {
		return view
	}

	scoreMsg := lipgloss.NewStyle().
		Width(50).
		Height(13).
		PaddingBottom(0).
		BorderForeground(charmtone.Blush).
		Foreground(charmtone.Blush).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render(strconv.Itoa(m.scoreA), " vs ", strconv.Itoa(m.scoreB))

	footer := lipgloss.NewStyle().
		Width(50).
		Height(13).
		PaddingBottom(0).
		BorderForeground(charmtone.Blush).
		Foreground(charmtone.Oyster).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render("Press Q to quit.")

	//:BUG: This could potentially lead to a bug because I am not sure if I did it correctly
	playGround := playGround("Something")
	finalUi := playGround + scoreMsg + footer

	centeredLayout := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		finalUi,
		// playGround,
	)

	comp := lipgloss.NewCompositor(
		lipgloss.NewLayer(playGround),
		lipgloss.NewLayer(scoreMsg),
		lipgloss.NewLayer(footer),
	)

	view.SetContent(comp.Render())
	return tea.NewView(centeredLayout)
}

func main() {
	// fmt.Println("Hello Mom")
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(69)
	}
}

//:References:
//https://github.com/charmbracelet/bubbletea/discussions/307
//https://github.com/charmbracelet/bubbletea/tree/main/examples/fullscreen
//https://golangdocs.com/concatenate-strings-in-golang
//https://www.google.com/search?q=setting+size+of+New+Compositor+bubble+tea+golang
