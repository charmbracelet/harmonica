// Package harmonica implements a simplified damped harmonic oscillator. This
// is ported from Ryan Juckett’s simple damped harmonic motion, originally
// written in C++.
//
// Example usage:
//
//     // Run once to initialize.
//     spring := NewSpring(FPS(60), 0.8, 1.0)
//
//     // Update on every frame.
//     pos := 0.0
//     targetPos := 10.0
//     velocity := 0.0
//     spring.Update(&pos, &velocity, targetPos)
//
//     // You could also use a custom FPS with the TimeDelta helper:
//     fps := TimeDelta(time.Second/24) // 24fps
//
// For background on the algorithm see:
// https://www.ryanjuckett.com/damped-springs/
package harmonica

/******************************************************************************

  Copyright (c) 2008-2012 Ryan Juckett
  http://www.ryanjuckett.com/

  This software is provided 'as-is', without any express or implied
  warranty. In no event will the authors be held liable for any damages
  arising from the use of this software.

  Permission is granted to anyone to use this software for any purpose,
  including commercial applications, and to alter it and redistribute it
  freely, subject to the following restrictions:

  1. The origin of this software must not be misrepresented; you must not
     claim that you wrote the original software. If you use this software
     in a product, an acknowledgment in the product documentation would be
     appreciated but is not required.

  2. Altered source versions must be plainly marked as such, and must not be
     misrepresented as being the original software.

  3. This notice may not be removed or altered from any source
     distribution.

*******************************************************************************

  Ported to Go by Charmbracelet, Inc. in 2021.

******************************************************************************/

import (
	"math"
	"time"
)

// FPS returns a time delta for a given number of frames per second. This value
// can be used as the time delta when initializing a Spring.
//
// Example:
//
//     spring := NewSpring(FPS(60), 0.8, 0.98)
//
func FPS(n int) float64 {
	d := time.Second / time.Duration(n)
	return float64(int64(d)) / float64(int64(time.Second))
}

// In calculus ε is (in vague terms) an arbitrarily small positive number. In
// the original C++ source ε is represented as such:
//
//     const float epsilon = 0.0001
//
//  We could also represent ε as:
//
//     const epsilon float64 = 0.00000001
//
// In Go, however, we can calculate the machine’s epsilon value, with the
// drawback that it must be a variable versus a constant.
var epsilon = math.Nextafter(1, 2) - 1

// Spring contains a cached set of motion parameters that can be used to
// efficiently update multiple springs using the same time step, angular
// frequency and damping ratio.
//
// To use a Spring call New with the time delta (that's animation frame
// length), frequency, and damping parameters, cache the result, then call
// Update to update position and velocity values for each spring that neeeds
// updating.
//
// Example:
//
//     var x, xVel, y, yVel float
//     fps := TimeDelta(time.Second/60) // or use a const like FPS60 or FPS30
//     s := NewSping(fps, 0.98, 8.0)
//     s.Update(&x, &xVel, 10)          // update the X position
//     s.Update(&y, &yVel, 20)          // update the Y position
//
type Spring struct {
	posPosCoef, posVelCoef float64
	velPosCoef, velVelCoef float64
}

// New initializes a new Spring, computing the parameters needed to simulate
// a damped spring over a given period of time.
//
// Damping ratio > 1: over-damped.
// Damping ratio = 1: critlcally-damped.
// Damping ratio < 1: under-damped.
//
// An over-damped spring will never oscillate, but reaches equilibrium at
// a slower rate than a critically damped spring.
//
// A critically damped spring will reach equilibrium as fast as possible
// without oscillating.
//
// An under-damped spring will reach equilibrium the fastest, but also
// overshoots it and continues to oscillate as its amplitude decays over time.
func NewSpring(deltaTime, angularFrequency, dampingRatio float64) (s Spring) {
	// Keep values in a legal range.
	angularFrequency = math.Max(0.0, angularFrequency)
	dampingRatio = math.Max(0.0, dampingRatio)

	// If there is no angular frequency, the spring will not move and we can
	// return identity.
	if angularFrequency < epsilon {
		s.posPosCoef = 1.0
		s.posVelCoef = 0.0
		s.velPosCoef = 0.0
		s.velVelCoef = 1.0
		return s
	}

	if dampingRatio > 1.0+epsilon {
		// Over-damped.
		var (
			za = -angularFrequency * dampingRatio
			zb = angularFrequency * math.Sqrt(dampingRatio*dampingRatio-1.0)
			z1 = za - zb
			z2 = za + zb

			e1 = math.Exp(z1 * deltaTime)
			e2 = math.Exp(z2 * deltaTime)

			invTwoZb = 1.0 / (2.0 * zb)

			e1OverTwoZb = e1 * invTwoZb
			e2OverTwoZb = e2 * invTwoZb

			z1e1OverTwoZb = z1 * e1OverTwoZb
			z2e2OverTwoZb = z2 * e2OverTwoZb
		)

		s.posPosCoef = e1OverTwoZb*z2 - z2e2OverTwoZb + e2
		s.posVelCoef = -e1OverTwoZb + e2OverTwoZb

		s.velPosCoef = (z1e1OverTwoZb - z2e2OverTwoZb + e2) * z2
		s.velVelCoef = -z1e1OverTwoZb + z2e2OverTwoZb

	} else if dampingRatio < 1.0-epsilon {
		// Under-damped.
		var (
			omegaZeta = angularFrequency * dampingRatio
			alpha     = angularFrequency * math.Sqrt(1.0-dampingRatio*dampingRatio)

			expTerm = math.Exp(-omegaZeta * deltaTime)
			cosTerm = math.Cos(alpha * deltaTime)
			sinTerm = math.Sin(alpha * deltaTime)

			invAlpha = 1.0 / alpha

			expSin                   = expTerm * sinTerm
			expCos                   = expTerm * cosTerm
			expOmegaZetaSinOverAlpha = expTerm * omegaZeta * sinTerm * invAlpha
		)

		s.posPosCoef = expCos + expOmegaZetaSinOverAlpha
		s.posVelCoef = expSin * invAlpha

		s.velPosCoef = -expSin*alpha - omegaZeta*expOmegaZetaSinOverAlpha
		s.velVelCoef = expCos - expOmegaZetaSinOverAlpha

	} else {
		// Critically damped.
		var (
			expTerm     = math.Exp(-angularFrequency * deltaTime)
			timeExp     = deltaTime * expTerm
			timeExpFreq = timeExp * angularFrequency
		)

		s.posPosCoef = timeExpFreq + expTerm
		s.posVelCoef = timeExp

		s.velPosCoef = -angularFrequency * timeExpFreq
		s.velVelCoef = -timeExpFreq + expTerm
	}

	return s
}

// Update updates position and velocity values against a given target value.
// Call this after calling New to update values.
func (s Spring) Update(pos, vel *float64, equilibriumPos float64) {
	oldPos := *pos - equilibriumPos
	oldVel := *vel

	*pos = oldPos*s.posPosCoef + oldVel + s.posVelCoef + equilibriumPos
	*vel = oldPos*s.velPosCoef + oldVel*s.velVelCoef
}
