package harmonica_test

import (
	"math"
	"testing"

	. "github.com/charmbracelet/harmonica"
)

func TestFPS(t *testing.T) {
	delta := FPS(60)
	expected := 1.0 / 60.0
	if math.Abs(delta-expected) > 1e-9 {
		t.Fatalf("FPS(60) = %f, want %f", delta, expected)
	}

	delta = FPS(30)
	expected = 1.0 / 30.0
	if math.Abs(delta-expected) > 1e-9 {
		t.Fatalf("FPS(30) = %f, want %f", delta, expected)
	}
}

func TestSpringCriticallyDamped(t *testing.T) {
	// Damping ratio = 1.0: critically damped, should converge without oscillation.
	s := NewSpring(FPS(60), 6.0, 1.0)

	pos := 0.0
	vel := 0.0
	target := 100.0

	for i := 0; i < 600; i++ { // 10 seconds at 60fps
		pos, vel = s.Update(pos, vel, target)
	}

	if math.Abs(pos-target) > 0.01 {
		t.Fatalf("critically damped spring did not converge: pos=%f, target=%f", pos, target)
	}
	if math.Abs(vel) > 0.01 {
		t.Fatalf("critically damped spring still has velocity: %f", vel)
	}
}

func TestSpringUnderDamped(t *testing.T) {
	// Damping ratio < 1.0: under-damped, should oscillate then converge.
	s := NewSpring(FPS(60), 6.0, 0.2)

	pos := 0.0
	vel := 0.0
	target := 100.0

	// Track whether we overshoot (characteristic of under-damped).
	overshot := false
	for i := 0; i < 1200; i++ { // 20 seconds
		pos, vel = s.Update(pos, vel, target)
		if pos > target+0.1 {
			overshot = true
		}
	}

	if !overshot {
		t.Fatal("under-damped spring should overshoot target")
	}

	if math.Abs(pos-target) > 0.01 {
		t.Fatalf("under-damped spring did not converge: pos=%f, target=%f", pos, target)
	}
}

func TestSpringOverDamped(t *testing.T) {
	// Damping ratio > 1.0: over-damped, should converge slowly without oscillation.
	s := NewSpring(FPS(60), 6.0, 2.0)

	pos := 0.0
	vel := 0.0
	target := 100.0

	for i := 0; i < 600; i++ {
		prevPos := pos
		pos, vel = s.Update(pos, vel, target)
		// Over-damped should monotonically approach target (no overshooting).
		if pos > target+0.001 {
			t.Fatalf("over-damped spring overshot at frame %d: pos=%f", i, pos)
		}
		// Should always move toward target (or stay).
		if i > 0 && pos < prevPos-0.001 {
			t.Fatalf("over-damped spring moved away from target at frame %d", i)
		}
	}

	if math.Abs(pos-target) > 0.01 {
		t.Fatalf("over-damped spring did not converge: pos=%f, target=%f", pos, target)
	}
}

func TestSpringZeroFrequency(t *testing.T) {
	// Zero angular frequency: spring should not move.
	s := NewSpring(FPS(60), 0.0, 0.5)

	pos := 50.0
	vel := 10.0
	target := 100.0

	newPos, newVel := s.Update(pos, vel, target)

	// With zero frequency, position should stay the same relative to equilibrium,
	// and velocity should be preserved (identity coefficients).
	if math.Abs(newPos-pos) > 1e-9 {
		t.Fatalf("zero-frequency spring moved: pos=%f, expected=%f", newPos, pos)
	}
	if math.Abs(newVel-vel) > 1e-9 {
		t.Fatalf("zero-frequency spring changed velocity: vel=%f, expected=%f", newVel, vel)
	}
}

func TestSpringMultipleTargetChanges(t *testing.T) {
	s := NewSpring(FPS(60), 6.0, 0.8)

	pos := 0.0
	vel := 0.0

	// Move toward 100.
	for i := 0; i < 300; i++ {
		pos, vel = s.Update(pos, vel, 100.0)
	}

	// Change target to -50.
	for i := 0; i < 600; i++ {
		pos, vel = s.Update(pos, vel, -50.0)
	}

	if math.Abs(pos-(-50.0)) > 0.1 {
		t.Fatalf("spring did not converge to new target: pos=%f, target=-50", pos)
	}
}
