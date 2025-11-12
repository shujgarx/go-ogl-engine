# Go OpenGL Mini-Engine (ECS + Physics + Audio)

A minimal cross‑platform starter engine in Go with:
- OpenGL 3.3 Core renderer (GLFW + go-gl)
- ECS (entities/components/systems)
- Simple 3D physics (gravity + floor collision)
- Audio (generated sine wave via Oto)
- Texture loading (PNG/JPEG)
- A feature-complete tech demo scene

## Quick start

```bash
# Requires a C compiler and OpenGL/GLFW dev headers
# macOS: xcode-select --install
# Ubuntu/Debian: sudo apt-get install -y libgl1-mesa-dev xorg-dev
# Windows: install MSYS2/LLVM; ensure gcc/clang in PATH

git init . && go mod download
go run ./cmd/game
```

## Tech demo showcase

The bundled sample scene now highlights the engine subsystems:

- Orbiting dynamic point light with per-pixel Phong shading.
- Animated "hero" cube hovering in the center to stress test the shader pipeline.
- Physically simulated stacks of cubes with gravity, damping, and bouncy floor collisions.
- Wide tiled floor mesh to demonstrate custom geometry and UV tiling.
- Automatic camera fly-around and ambient sine tone playback.

## Windows portable build

To assemble a distributable build (executable + assets + DLLs) on Windows, run:

```
build_techdemo.cmd
```

The script will:

- Ensure the Go toolchain and a `gcc`/MinGW-w64 compiler are available.
- Download the official GLFW prebuilt package automatically on first run.
- Fetch Go module dependencies, compile the tech demo, and copy assets/DLLs into `Release\`.
- Produce `Release\bin\TechDemo.exe` ready to ship alongside the assets folder.

### Notes
- We lock the main goroutine to the OS thread as OpenGL contexts are thread-bound.
- On macOS, OpenGL is deprecated; we target 3.3 Core and set ForwardCompatible = true.
