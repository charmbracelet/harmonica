package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

const (
	width       = 1024
	height      = 768
	bgColor     = "#575BD8"
	textColor   = "#827EFF"
	spriteColor = "#FFFDF5"
	fontFile    = "JetBrainsMono-Regular.ttf"
	idleTimeout = time.Second * 4
)

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(Game{}.Run)
}

type Game struct {
	deltaTime float64
	frequency float64
	damping   float64
	win       *pixelgl.Window
	shift     bool
	font      font.Face
	sprite    *Sprite
	lastClick time.Time
	dirty     bool
	showHelp  bool
}

func (g Game) Size() (width, height float64) {
	return g.win.Bounds().W(), g.win.Bounds().H()
}

func (g Game) MousePosition() (x, y float64) {
	return g.win.MousePosition().X, g.win.MousePosition().Y
}

func (g Game) Run() {
	g.frequency = 10.0
	g.damping = 0.2
	g.showHelp = true

	var err error
	if g.font, err = gg.LoadFontFace(fontFile, 24); err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:     "Harmonica Example",
		Bounds:    pixel.R(0, 0, width, height),
		VSync:     true,
		Resizable: false,
	}

	if g.win, err = pixelgl.NewWindow(cfg); err != nil {
		panic(err)
	}

	last := time.Now()
	for !g.win.Closed() {
		w, h := g.Size()

		// Get delta time
		g.deltaTime = time.Since(last).Seconds()
		last = time.Now()

		g.Update()

		// Use fogleman/gg to render everything
		ctx := gg.NewContext(int(w), int(h))
		g.Draw(ctx)
		canvas := g.win.Canvas()
		canvas.SetPixels(ctx.Image().(*image.RGBA).Pix)

		// Render
		g.win.Update()
	}
}

func (g *Game) Update() {
	pressed := g.win.JustPressed
	released := g.win.JustReleased

	// Handle shift key
	switch {
	case pressed(pixelgl.KeyLeftShift), pressed(pixelgl.KeyRightShift):
		g.shift = true
	case released(pixelgl.KeyLeftShift), released(pixelgl.KeyRightShift):
		g.shift = false
	}

	adjustFreq, adjustDamp := 0.0, 0.0

	// Handle arrow keys to adjust frequency and damping
	switch {
	case released(pixelgl.KeyLeft):
		adjustFreq = -0.1
	case released(pixelgl.KeyRight):
		adjustFreq = 0.1
	case released(pixelgl.KeyUp):
		adjustDamp = 0.01
	case released(pixelgl.KeyDown):
		adjustDamp = -0.01
	case released(pixelgl.KeySpace):
		g.showHelp = !g.showHelp
	}

	if g.shift {
		adjustDamp *= 10
		adjustFreq *= 10
	}

	if adjustFreq != 0 || adjustDamp != 0 {
		g.frequency = math.Max(0, g.frequency+adjustFreq)
		g.damping = math.Max(0, g.damping+adjustDamp)
		g.dirty = true
	}

	if time.Now().After(g.lastClick.Add(idleTimeout)) {
		if g.sprite != nil {
			g.sprite.Hide()
		}
	}

	if released(pixelgl.MouseButtonLeft) {
		if g.sprite == nil {
			g.sprite = NewSprite(g)
		} else {
			s := g.sprite
			s.TargetX, s.TargetY = g.MousePosition()

			switch s.State() {
			case hiding, gone:
				s.Show()
			}
		}

		g.lastClick = time.Now()
	}

	// Sprite
	if g.sprite != nil {
		g.sprite.Update()
	}

	g.dirty = false
}

func (g Game) Draw(ctx *gg.Context) {
	w, h := g.Size()

	// BG
	ctx.SetHexColor(bgColor)
	ctx.DrawRectangle(0, 0, w, h)
	ctx.Fill()

	// Help text
	if g.showHelp {
		// For some reason our text renders upside down and backwards. Apply a
		// matrix to fix it.
		ctx.ScaleAbout(-1, 1, w/2, h/2)
		ctx.RotateAbout(math.Pi, w/2, h/2)

		// Instructional text
		ctx.Push()
		ctx.SetHexColor(textColor)
		ctx.SetFontFace(g.font)
		ctx.DrawStringAnchored("Click!", w/2, h/2, 0.5, 0.5)
		ctx.Fill()
		ctx.Pop()

		// Status
		str := fmt.Sprintf(
			"Frequency: %.1f (←/→: adjust) • Damping: %.2f (↑/↓: adjust)",
			g.frequency, g.damping,
		)
		ctx.Push()
		ctx.SetHexColor(textColor)
		ctx.SetFontFace(g.font)
		ctx.DrawStringAnchored(str, w/2, h-34, 0.5, 0)
		ctx.Fill()
		ctx.Pop()

		// Reset matrix
		ctx.Identity()
	}

	// Draw sprites
	if g.sprite != nil {
		g.sprite.Draw(ctx)
	}
}

type SpriteState int

const (
	moving SpriteState = iota
	stopped
	hiding
	gone
)

type Sprite struct {
	TargetX, TargetY, TargetRadius float64
	X, xVel                        float64
	Y, yVel                        float64
	radius, radiusVel              float64
	color                          string
	spring                         harmonica.Spring
	game                           *Game
}

func NewSprite(g *Game) *Sprite {
	mouseX, mouseY := g.MousePosition()

	s := &Sprite{
		X:       mouseX,
		Y:       mouseY,
		TargetX: mouseX,
		TargetY: mouseY,
		radius:  0.1,
		color:   spriteColor,
		game:    g,
	}
	s.computeSpring()
	s.randomRadius()

	return s
}

func (s *Sprite) Update() {
	if s.State() == gone {
		return
	}

	if s.game.dirty {
		// Recompute spring coefficients since our frequency or damping has
		// changed.
		s.computeSpring()
	}

	// Calculate positions based on our spring
	s.X, s.xVel = s.spring.Update(s.X, s.xVel, s.TargetX)
	s.Y, s.yVel = s.spring.Update(s.Y, s.yVel, s.TargetY)
	s.radius, s.radiusVel = s.spring.Update(s.radius, s.radiusVel, s.TargetRadius)
}

func (s Sprite) Draw(ctx *gg.Context) {
	if s.State() == gone {
		return
	}

	ctx.Push()
	ctx.DrawCircle(s.X, s.Y, s.radius)
	ctx.SetHexColor(s.color)
	ctx.Fill()
	ctx.Pop()
}

func (s Sprite) State() SpriteState {
	const precision = 2
	x := roundFloat(s.X, precision)
	tX := roundFloat(s.TargetX, precision)
	y := roundFloat(s.Y, precision)
	tY := roundFloat(s.TargetY, precision)
	r := roundFloat(s.radius, precision)
	tR := roundFloat(s.TargetRadius, precision)

	if r == 0.0 {
		return gone
	}
	if s.TargetRadius == 0 {
		return hiding
	}
	if x == tX && y == tY && r == tR {
		return stopped
	}
	return moving
}

func (s *Sprite) computeSpring() {
	// Calculate spring coefficients
	s.spring = harmonica.NewSpring(s.game.deltaTime, s.game.frequency, s.game.damping)
}

func (s *Sprite) randomRadius() {
	s.TargetRadius = rand.Float64()*100.0 + 100.0
}

func (s *Sprite) Show() {
	if s.State() == gone {
		s.X, s.Y = s.game.MousePosition()
	}
	if s.radius < 0.1 {
		s.radius = 0.1
	}
	s.computeSpring()
	s.randomRadius()
}

func (s *Sprite) Hide() {
	// Fixed frequency and damping when hiding
	s.spring = harmonica.NewSpring(s.game.deltaTime, 32.0, 1.0)
	s.TargetRadius = 0
}

func roundFloat(input float64, decimalPlaces int) float64 {
	pow := math.Pow(10, float64(decimalPlaces))
	return math.Round(pow*input) / pow
}
