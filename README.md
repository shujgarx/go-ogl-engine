# Go OpenGL Mini-Engine (ECS + Physics + Audio)

A minimal cross‑platform starter engine in Go with:
- OpenGL 3.3 Core renderer (GLFW + go-gl)
- ECS (entities/components/systems)
- Simple 3D physics (gravity + floor collision)
- Audio (generated sine wave via Oto)
- Texture loading (PNG/JPEG)
- A rotating textured cube demo

## Quick start

```bash
# Requires a C compiler and OpenGL/GLFW dev headers
# macOS: xcode-select --install
# Ubuntu/Debian: sudo apt-get install -y libgl1-mesa-dev xorg-dev
# Windows: install MSYS2/LLVM; ensure gcc/clang in PATH

git init . && go mod download
go run ./cmd/game
```

### Notes
- We lock the main goroutine to the OS thread as OpenGL contexts are thread-bound.
- On macOS, OpenGL is deprecated; we target 3.3 Core and set ForwardCompatible = true.
