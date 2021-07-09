package main

import (
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/h8gi/canvas"
)

const (
	canvasWidth  = 1024
	canvasHeight = 768
	fps          = 60
	circleRadius = 100
	bgColor      = "#111111"
	idleTimeout  = time.Second * 5
)

type State int

const (
	unstarted State = iota
	running
	timeout
)

func main() {
	var (
		circ      Circle
		state     State
		lastClick time.Time
	)

	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     canvasWidth,
		Height:    canvasHeight,
		FrameRate: fps,
		Title:     "Harmonica",
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.Clear()
		ctx.SetHexColor(bgColor)

		switch state {
		case unstarted:
			if ctx.IsMouseDragged {
				const frequency = 0.8
				const damping = 1.0

				circ = Circle{
					// Setup a new spring.
					spring: harmonica.NewSpring(harmonica.TimeDelta(time.Second/fps), frequency, damping),

					// Set target radius.
					Radius: circleRadius,
				}

				// Set starting position.
				circ.SetPos(ctx.Mouse.X, ctx.Mouse.Y)

				lastClick = time.Now()
				state = running
			}

		case running:
			if ctx.IsMouseDragged {
				lastClick = time.Now()

				if circ.Hidden {
					circ.SetPos(ctx.Mouse.X, ctx.Mouse.Y)
					circ.Hidden = false
				}

				// Update targets for our spring animation.
				circ.X = ctx.Mouse.X
				circ.Y = ctx.Mouse.Y
				circ.Radius = circleRadius
			}

			circ.Draw(ctx)

			if time.Now().After(lastClick.Add(idleTimeout)) {
				state = timeout
			}

		case timeout:
			circ.Radius = 0
			circ.Draw(ctx)
			if circ.rad == 0 {
				state = running
				lastClick = time.Now()
			}
		}
	})
}

type Circle struct {
	// Target values. This is what the spring will use as an equilibrium to
	// animate towards.
	X, Y, Radius float64

	x, xVel     float64
	y, yVel     float64
	rad, radVel float64
	spring      harmonica.Spring

	Hidden bool
}

func (c *Circle) SetPos(x, y float64) {
	c.x = x
	c.y = y
	c.Hidden = false
}

func (c *Circle) Draw(ctx *canvas.Context) {
	// Update position and radius.
	c.spring.Update(&c.x, &c.xVel, c.X)
	c.spring.Update(&c.y, &c.yVel, c.Y)
	c.spring.Update(&c.rad, &c.radVel, c.Radius)

	if c.rad < 0 {
		c.rad = 0
		c.radVel = 0
		c.Hidden = true
	}

	const color = "#f1f1f1"

	ctx.Push()
	ctx.DrawCircle(c.x, c.y, c.rad)
	ctx.SetHexColor(color)
	ctx.Fill()
	ctx.Pop()
}
