package harmonica_test

import (
	"math"
	"testing"

	. "github.com/charmbracelet/harmonica"
)

const fps = 60

func TestNew(t *testing.T) {
	x := 8
	y := 20
	z := 0

	projectile := NewProjectile(FPS(60), Point{float64(x), float64(y), float64(z)}, Vector{1, 1, 0}, Vector{0, 9.81, 0})
	pos := projectile.Update()
	if x != int(pos.X) {
		t.Logf("Want: %d, Got: %d", int(x), int(pos.X))
		t.Fatal("x coordinate unexpected")
	}

	if y != int(pos.Y) {
		t.Logf("Want: %d, Got: %d", int(y), int(pos.Y))
		t.Fatal("y coordinate unexpected")
	}
}

const equalityThreshold = 1e-2

// floating point comparison function that tests for an equality under the:
//    equalityThreshold
func equal(a, b float64) bool {
	return math.Abs(a-b) <= equalityThreshold
}

func TestUpdate(t *testing.T) {
	projectile := NewProjectile(FPS(fps), Point{0, 0, 0}, Vector{5, 5, 0}, Vector{0, 0, 0})
	coordinates := []Point{
		{5.0, 5.0, 0},
		{10.0, 10.0, 0},
		{15.0, 15.0, 0},
		{20.0, 20.0, 0},
		{25.0, 25.0, 0},
		{30.0, 30.0, 0},
		{35.0, 35.0, 0},
	}

	for _, c := range coordinates {
		var pos Point
		for i := 0; i < fps; i++ {
			pos = projectile.Update()
		}

		pvel := projectile.Velocity()
		if !equal(pvel.X, 5) || !equal(pvel.Y, 5) || !equal(pvel.Z, 0) {
			t.Logf("Want: (%.2f, %.2f, %.2f)", pvel.X, pvel.Y, pvel.Z)
			t.Logf("Want: (%.2f, %.2f, %.2f)", 5.0, 5.0, 0.0)
			t.Fatal("velocity unexpected")
		}

		if !equal(pos.X, c.X) || !equal(pos.Y, c.Y) {
			t.Logf("Want: (%.2f, %.2f)", c.X, c.Y)
			t.Logf("Got:  (%.2f, %.2f)", pos.X, pos.Y)
			t.Fatal("coordinate unexpected")
		}
	}
}

func TestUpdateGravity(t *testing.T) {
	fps := 60
	projectile := NewProjectile(FPS(fps), Point{0, 0, 0}, Vector{5, 5, 0}, TerminalGravity)

	coordinates := []Point{
		{5.0, 9.82, 0},
		{10.0, 29.46, 0},
		{15.0, 58.90, 0},
		{20.0, 98.15, 0},
		{25.0, 147.22, 0},
		{30.0, 206.09, 0},
		{35.0, 274.77, 0},
	}

	for _, c := range coordinates {
		var pos Point
		for f := 0; f < fps; f++ {
			pos = projectile.Update()
		}

		yacc := projectile.Acceleration().Y
		if !equal(yacc, TerminalGravity.Y) {
			t.Logf("Want: %.2f", TerminalGravity.Y)
			t.Logf("Got:  %.2f", yacc)
			t.Fatal("Y acceleration unexpected")
		}

		if !equal(pos.X, c.X) || !equal(pos.Y, c.Y) {
			t.Logf("Want: (%.2f, %.2f)", c.X, c.Y)
			t.Logf("Got:  (%.2f, %.2f)", pos.X, pos.Y)
			t.Fatal("coordinate unexpected")
		}
	}
}
