package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps       = 60
	maxHeight = 100
)

var spriteStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#575BD8"))

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
	projectile *harmonica.Projectile
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

	out.WriteString(strings.Repeat("\n", int(m.pos.Y)))

	out.WriteString(strings.Repeat(" ", int(m.pos.X)))
	out.WriteString(spriteStyle.Render("//////////") + "\n")
	out.WriteString(strings.Repeat(" ", int(m.pos.X)))
	out.WriteString(spriteStyle.Render("//////////") + "\n")
	out.WriteString(strings.Repeat(" ", int(m.pos.X)))
	out.WriteString(spriteStyle.Render("//////////") + "\n")
	out.WriteString(strings.Repeat(" ", int(m.pos.X)))
	out.WriteString(spriteStyle.Render("//////////") + "\n")
	out.WriteString(strings.Repeat(" ", int(m.pos.X)))
	out.WriteString(spriteStyle.Render("//////////") + "\n")

	return out.String()
}

func main() {
	initPos := harmonica.Point{X: 0, Y: 0}
	initVel := harmonica.Vector{X: 0, Y: 0}
	initAcc := harmonica.Vector{X: 250, Y: 0}
	m := model{
		projectile: harmonica.NewProjectile(harmonica.FPS(fps), initPos, initVel, initAcc),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
