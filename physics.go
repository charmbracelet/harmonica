// Simple physics projectile motion.
//
// Example usage:
//
//    // Run once to initialize.
//    projectile := NewProjectile(FPS(60), Point{6.0, 100.0, 0.0}, Vector{2.0, 0.0, 0.0}, Vector{2.0, -9.81, 0.0})
//
//    // Update on every frame.
//    someUpdateLoop(func() {
//      pos := projectile.Update()
//    })
//
// For background on projectile motion see:
// https://en.wikipedia.org/wiki/Projectile_motion
package harmonica

// Projectile is the representation of a projectile that has a position on a
// plane and an acceleration and velocity
type Projectile struct {
	pos       Point
	vel       Vector
	acc       Vector
	deltaTime float64
}

// Point is a representation of a point which contains the X, Y, Z coordinates
// of the point on a plane.
type Point struct {
	X, Y, Z float64
}

// Vector is a representation of a vector which carries a magnitude and a
// direction.  We represent the vector as a point from the origin (0, 0) where
// the magnitude is the euclidean distance from the origin and the direction is
// the direction to the point from the origin.
type Vector struct {
	X, Y, Z float64
}

// Gravity is a utility vector that represents gravity in 2D and 3D contexts,
// assuming that your coordinate plane looks like in 2D or 3D:
//
//  -y            -y ±z
//   │             │ /
//   │             │/
//   └───── ±x     └───── ±x
//
// Note: Gravity usually is -9.81m/s however we use a positive value because we
// assume the origin is placed at the top-left corner and that downward is the
// positive y direction (only if using this utility variable).
// Otherwise, you can place the origin wherever you'd like.
var Gravity = Vector{0, 9.81, 0}

// NewProjectile accepts a frame rate, and initial values for position, velocity, and acceleration and
// returns a new projectile.
func NewProjectile(deltaTime float64, initialPosition Point, initialVelocity, initalAcceleration Vector) Projectile {
	return Projectile{
		pos:       initialPosition,
		vel:       initialVelocity,
		acc:       initalAcceleration,
		deltaTime: deltaTime,
	}
}

// Update updates the position and velocity values for the given projectile.
// Call this after calling NewProjectile to update values.
func (p *Projectile) Update() Point {
	p.pos.X += (p.vel.X * p.deltaTime)
	p.pos.Y += (p.vel.Y * p.deltaTime)
	p.pos.Z += (p.vel.Z * p.deltaTime)

	p.vel.X += (p.acc.X * p.deltaTime)
	p.vel.Y += (p.acc.Y * p.deltaTime)
	p.vel.Z += (p.acc.Z * p.deltaTime)

	return p.pos
}
