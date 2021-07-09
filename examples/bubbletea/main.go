package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps          = time.Second / 60
	spriteWidth  = 12
	spriteHeight = 5
	frequency    = 0.95
	damping      = 0.98
)

var (
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "246", Dark: "241"}).
			MarginTop(1).
			MarginLeft(2)

	spriteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("69"))
)

type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(fps, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

type model struct {
	x      float64
	xVel   float64
	spring harmonica.Spring
}

func (_ model) Init() tea.Cmd {
	return animate()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	// Step foreward one frame
	case frameMsg:
		const targetX = 60

		// Update x position (and velocity) with our spring.
		m.spring.Update(&m.x, &m.xVel, targetX)

		// Quit when we're basically at the target position.
		if math.Abs(m.x-targetX) < 0.01 {
			return m, tea.Quit
		}

		// Request next frame
		return m, animate()

	default:
		return m, nil
	}
}

func (m model) View() string {
	var out strings.Builder
	fmt.Fprint(&out, "\n")

	x := int(math.Round(m.x))
	if x < 0 {
		return ""
	}

	spriteRow := spriteStyle.Render(strings.Repeat("/", spriteWidth))
	row := strings.Repeat(" ", x) + spriteRow + "\n"
	fmt.Fprint(&out, strings.Repeat(row, spriteHeight))

	fmt.Fprint(&out, helpStyle.Render("Press any key to quit"))

	return out.String()
}

func main() {
	m := model{
		spring: harmonica.NewSpring(harmonica.TimeDelta(fps), frequency, damping),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
