package main

import (
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/h8gi/canvas"
)

const (
	width   = 1024
	height  = 768
	fps     = 60
	bgColor = "#111111"
)

func main() {
	var (
		circ       = newCircle()
		clickedYet bool
	)

	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     width,
		Height:    height,
		FrameRate: fps,
		Title:     "Harmonica",
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.SetHexColor(bgColor)
		ctx.Clear()

		if ctx.IsMouseDragged {
			if !clickedYet {
				// Set starting positions.
				circ.x = ctx.Mouse.X
				circ.y = ctx.Mouse.Y

				circ.Radius = 100
				clickedYet = true
			}
			circ.X = ctx.Mouse.X
			circ.Y = ctx.Mouse.Y
		}

		if clickedYet {
			circ.draw(ctx)
		}
	})
}

type circle struct {
	// Target values. This is what the spring will seek.
	X, Y, Radius float64

	// Managed values. These are the initial
	x, xVel     float64
	y, yVel     float64
	rad, radVel float64

	spring harmonica.Spring
}

func newCircle() circle {
	const (
		frequency = 0.8
		damping   = 0.99
	)

	return circle{
		spring: harmonica.NewSpring(harmonica.TimeDelta(time.Second/fps), frequency, damping),
	}
}

func (c *circle) draw(ctx *canvas.Context) {
	// Update
	c.spring.Update(&c.x, &c.xVel, c.X)
	c.spring.Update(&c.y, &c.yVel, c.Y)
	c.spring.Update(&c.rad, &c.radVel, c.Radius)

	const color = "#f1f1f1"

	ctx.Push()
	ctx.DrawCircle(c.x, c.y, c.rad)
	ctx.SetHexColor(color)
	ctx.Fill()
}
