Harmonica
=========

<p>
    <a href="https://github.com/charmbracelet/harmonica/releases"><img src="https://img.shields.io/github/release/charmbracelet/harmonica.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/harmonica?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/charmbracelet/harmonica/actions"><img src="https://github.com/charmbracelet/harmonica/workflows/build/badge.svg" alt="Build Status"></a>
</p>

A simple, efficient spring animation library for smooth, natural motion.

```go
import "github.com/charmbracelet/harmonica"

// A thing we want to animate.
sprite := struct{
    x, xVelocity float64
    y, yVelocity float64
}{}

// Where we want to animate it.
const targetX = 50.0
const targetY = 100.0

// Initialize a spring with framerate, angular frequency, and damping values.
spring := harmonica.NewSpring(harmonica.FPS60, 0.8, 0.98)

// Animate!
for {
    spring.Update(&sprite.x, &sprite.xVelocity, targetX)
    spring.Update(&sprite.y, &sprite.yVelocity, targetY)
    time.Sleep(time.Second/60)
}
```

For details, see the [examples][examples] and the [docs][docs].

[examples]: https://github.com/charmbracelet/harmonica/tree/master/examples
[docs]: https://pkg.go.dev/github.com/charmbracelet/harmonica?tab=doc

## Acknowledgements

This library is a fairly straightforward port of [Ryan Juckett][juckett]’s
excellent damped simple harmonic oscillator originally writen in C++ in 2008
and published in 2012. [Ryan’s writeup][writeup] on the subject is fantastic.

[juckett]: https://www.ryanjuckett.com/
[writeup]: https://www.ryanjuckett.com/damped-springs/

## License

[MIT](https://github.com/charmbracelet/harmonica/raw/master/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge-unrounded.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
