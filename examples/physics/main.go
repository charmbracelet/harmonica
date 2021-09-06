package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
)

const (
	fps       = 60
	maxHeight = 100
)

type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

func wait(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return nil
	}
}

type model struct {
	projectile harmonica.Projectile
	pos        harmonica.Point
}

func (_ model) Init() tea.Cmd {
	return tea.Sequentially(wait(time.Second/2), animate())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	// Step forward one frame
	case frameMsg:
		m.pos = m.projectile.Update()

		if m.pos.Y > maxHeight {
			return m, tea.Quit
		}

		return m, animate()
	default:
		return m, nil
	}
}

func (m model) View() string {
	var out strings.Builder

	for y := 0; y < int(m.pos.Y); y++ {
		out.WriteString("\n")
	}

	for x := 0; x < int(m.pos.X); x++ {
		out.WriteString(" ")
	}
	out.WriteString(fmt.Sprintf("(%.2f, %.2f)", m.pos.X, m.pos.Y))

	return out.String()
}

func main() {
	initPos := harmonica.Point{X: 0, Y: 0}
	initVel := harmonica.Vector{X: 5, Y: 0}
	initAcc := harmonica.TerminalGravity
	m := model{
		projectile: harmonica.NewProjectile(harmonica.FPS(fps), initPos, initVel, initAcc),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
