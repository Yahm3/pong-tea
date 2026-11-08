package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

const (
	playW      = 80
	playH      = 40
	innerW     = playW - 2
	innerH     = playH - 2
	paddleH    = 4
	paddleStep = 2
)

const (
	startDelay = 80 * time.Millisecond
	minDelay   = 25 * time.Millisecond
	speedUp    = 4 * time.Millisecond
)

type tickMsg time.Time

type model struct {
	width, height    int
	quiting          bool
	ballX, ballY     int
	velX, velY       int
	paddleA, paddleB int
	scoreA, scoreB   int
	delay            time.Duration
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func newModel() model {
	m := model{
		paddleA: (innerH - paddleH) / 2,
		paddleB: (innerH - paddleH) / 2,
	}
	return m.serve(1)
}

func (m model) serve(dir int) model {
	m.ballX, m.ballY = innerW/2, innerH/2
	m.velX, m.velY = dir, 1
	m.delay = startDelay
	return m
}

func (m model) Init() tea.Cmd { return tick(startDelay) }

func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }

func hit(ballY, paddleTop int) bool {
	return ballY >= paddleTop && ballY < paddleTop+paddleH
}

func (m model) step() model {
	m.ballX += m.velX
	m.ballY += m.velY

	if m.ballY <= 0 {
		m.ballY, m.velY = 0, 1
	}

	if m.ballY >= innerH-1 {
		m.ballY, m.velY = innerH-1, -1
	}

	switch {
	case m.velX < 0 && m.ballX <= 1 && hit(m.ballY, m.paddleA):
		m.ballX, m.velX = 1, 1
		m.delay = max(minDelay, m.delay-speedUp)

	case m.velX > 0 && m.ballX >= innerW-2 && hit(m.ballY, m.paddleB):
		m.ballX, m.velX = innerW-2, -1
		m.delay = max(minDelay, m.delay-speedUp)

	case m.ballX < 0:
		m.scoreB++
		return m.serve(-1)

	case m.ballX >= innerW:
		m.scoreA++
		return m.serve(1)
	}

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tickMsg:
		return m.step(), tick(m.delay)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quiting = true
			return m, tea.Quit
		case "w":
			m.paddleA = clamp(m.paddleA-paddleStep, 0, innerH-paddleH)
		case "s":
			m.paddleA = clamp(m.paddleA+paddleStep, 0, innerH-paddleH)
		case "o":
			m.paddleB = clamp(m.paddleB-paddleStep, 0, innerH-paddleH)
		case "k":
			m.paddleB = clamp(m.paddleB+paddleStep, 0, innerH-paddleH)
		}
	}
	return m, nil
}

var (
	playgroundStyle = lipgloss.NewStyle().
			Width(playW).
			Height(playH).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(charmtone.Charple).
			Align(lipgloss.Center, lipgloss.Center)

	scoreStyle = lipgloss.NewStyle().
			Width(playW).
			Foreground(charmtone.Blush).
			Align(lipgloss.Center)

	footerStyle = lipgloss.NewStyle().
			Width(playW).
			Foreground(charmtone.Oyster).
			Align(lipgloss.Center)

	paddleStyle = lipgloss.NewStyle().Foreground(charmtone.Blush)
)

func (m model) View() tea.View {
	var view tea.View

	view.AltScreen = true

	if m.quiting {
		return view
	}

	score := scoreStyle.Render(fmt.Sprintf("%d vs %d", m.scoreA, m.scoreB))
	playground := playgroundStyle.Render("")
	footer := footerStyle.Render("W/S: left O/K: right Q: quit")

	pong := lipgloss.JoinVertical(lipgloss.Center, score, playground, footer)

	originX := 1
	originY := lipgloss.Height(score) + 1

	paddle := paddleStyle.Render(strings.TrimSuffix(strings.Repeat("#\n", paddleH), "\n"))

	left := lipgloss.NewLayer(paddle).X(originX).Y(originY + m.paddleA).Z(1)
	right := lipgloss.NewLayer(paddle).X(originX + innerW - 1).Y(originY + m.paddleB).Z(1)
	ball := lipgloss.NewLayer("O").X(originX + m.ballX).Y(originY + m.ballY).Z(2)

	comp := lipgloss.NewCompositor(
		lipgloss.NewLayer(pong), left, right, ball,
	)

	view.SetContent(lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		comp.Render(),
	))
	return view
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(69)
	}
}

//:References:
//https://github.com/charmbracelet/bubbletea/discussions/307
//https://github.com/charmbracelet/bubbletea/tree/main/examples/fullscreen
//https://golangdocs.com/concatenate-strings-in-golang
//https://www.google.com/search?q=setting+size+of+New+Compositor+bubble+tea+golang
