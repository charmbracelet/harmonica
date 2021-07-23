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
	fps          = 60
	spriteWidth  = 12
	spriteHeight = 5
	frequency    = 7.0
	damping      = 0.15
)

var (
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "246", Dark: "241"}).
			MarginTop(1).
			MarginLeft(2)

	spriteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#575BD8"))
)

type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

func waitASec(ms int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Millisecond * time.Duration(ms))
		return nil
	}
}

type model struct {
	x      float64
	xVel   float64
	spring harmonica.Spring
}

func (_ model) Init() tea.Cmd {
	return tea.Sequentially(waitASec(500), animate())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	// Step foreward one frame
	case frameMsg:
		const targetX = 60

		// Update x position (and velocity) with our spring.
		m.x, m.xVel = m.spring.Update(m.x, m.xVel, targetX)

		// Quit when we're basically at the target position.
		if math.Abs(m.x-targetX) < 0.01 {
			return m, tea.Sequentially(waitASec(750), tea.Quit)
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
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
