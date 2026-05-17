# Examples

Run examples from this directory (`examples/`), not the repository root.

## Terminal UI

```bash
go run ./spring/tui
go run ./particle
```

## OpenGL (`spring/opengl`)

This demo uses CGO and links against GLFW/OpenGL via `pkg-config`. Install system libraries first:

| OS | Packages |
|----|----------|
| macOS | `brew install glfw pkg-config` |
| Debian/Ubuntu | `sudo apt install libglfw3-dev libgl1-mesa-dev pkg-config` |
| Fedora | `sudo dnf install glfw-devel mesa-libGL-devel pkg-config` |

Then:

```bash
go run ./spring/opengl
```

If `pkg-config` cannot find `gl`, the build fails before Go compiles the example — that usually means the packages above are missing or `PKG_CONFIG_PATH` is unset.
