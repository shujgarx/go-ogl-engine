package engine

import (
	"runtime"
	"time"

	"example.com/go-ogl-engine/internal/audio"
	"example.com/go-ogl-engine/internal/ecs"
	"example.com/go-ogl-engine/internal/gfx"
	"example.com/go-ogl-engine/internal/physics"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	Window   *glfw.Window
	Renderer *gfx.Renderer
	Assets   string

	World      *ecs.World
	Transforms *ecs.Store[ecs.Transform]
	Renderers  *ecs.Store[ecs.MeshRenderer]
	Bodies     *ecs.Store[ecs.RigidBody]

	Phys  *physics.World
	Audio *audio.System

	MeshCube *gfx.Mesh
	Tex      *gfx.Texture
}

func New(width, height int, title, assets string) (*Engine, error) {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return nil, err
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}
	gl.Enable(gl.DEPTH_TEST)

	r := &Engine{Window: win, Assets: assets}
	r.Renderer, err = gfx.NewRenderer(assets, width, height)
	if err != nil {
		return nil, err
	}

	// ECS stores
	w := ecs.NewWorld()
	r.World = w
	r.Transforms = ecs.NewStore[ecs.Transform]()
	r.Renderers = ecs.NewStore[ecs.MeshRenderer]()
	r.Bodies = ecs.NewStore[ecs.RigidBody]()

	// Physics
	r.Phys = physics.New(r.Bodies, r.Transforms)

	// Audio
	r.Audio = audio.NewSystem()
	_ = r.Audio.Init(48000)

	// Content
	r.MeshCube = makeCube()
	r.Tex, _ = gfx.LoadTexture(assets + "/textures/checker.png")

	// Scene
	cube := w.NewEntity()
	r.Transforms.Set(cube, ecs.Transform{Position: [3]float32{0, 0, 0}, Rotation: [3]float32{0, 0, 0}, Scale: [3]float32{1, 1, 1}})
	r.Renderers.Set(cube, ecs.MeshRenderer{MeshID: "cube", TextureID: "checker"})
	r.Bodies.Set(cube, ecs.RigidBody{UseGravity: true, Mass: 1.0})

	// Floor as invisible collider
	floor := w.NewEntity()
	r.Transforms.Set(floor, ecs.Transform{Position: [3]float32{0, -1, 0}, Scale: [3]float32{5, 0.1, 5}})
	return r, nil
}

func (e *Engine) Run() {
	last := time.Now()
	angle := float32(0)
	_ = e.Audio.PlaySine(440, 0.2, 0.1)

	for !e.Window.ShouldClose() {
		now := time.Now()
		dt := float32(now.Sub(last).Seconds())
		last = now
		angle += dt

		// Update physics
		e.Phys.Step(dt)

		// Spin the cube
		e.Transforms.ForEach(func(ent ecs.Entity, t *ecs.Transform) {
			if e.Renderers.Has(ent) {
				t.Rotation[1] = angle
			}
		})

		// Render
		gl.ClearColor(0.1, 0.12, 0.15, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(e.Renderer.Program)
		// Bind texture 0
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, e.Tex.ID)

		e.Transforms.ForEach(func(ent ecs.Entity, t *ecs.Transform) {
			if mr, ok := e.Renderers.Get(ent); ok {
				_ = mr // single mesh
				mvp := e.Renderer.Cam.MVP(t.Position, t.Rotation, t.Scale)
				gl.UniformMatrix4fv(e.Renderer.UniformMVP, 1, false, &mvp[0])
				e.MeshCube.Draw()
			}
		})

		e.Window.SwapBuffers()
		glfw.PollEvents()
	}
	e.Shutdown()
}

func (e *Engine) Shutdown() {
	glfw.Terminate()
}

// Hard-coded cube (pos, normal, uv)
func makeCube() *gfx.Mesh {
	v := []float32{
		// positions        // normals         // uvs
		// front
		-0.5, -0.5, 0.5, 0, 0, 1, 0, 0,
		0.5, -0.5, 0.5, 0, 0, 1, 1, 0,
		0.5, 0.5, 0.5, 0, 0, 1, 1, 1,
		-0.5, 0.5, 0.5, 0, 0, 1, 0, 1,
		// back
		-0.5, -0.5, -0.5, 0, 0, -1, 1, 0,
		0.5, -0.5, -0.5, 0, 0, -1, 0, 0,
		0.5, 0.5, -0.5, 0, 0, -1, 0, 1,
		-0.5, 0.5, -0.5, 0, 0, -1, 1, 1,
		// left
		-0.5, -0.5, -0.5, -1, 0, 0, 0, 0,
		-0.5, -0.5, 0.5, -1, 0, 0, 1, 0,
		-0.5, 0.5, 0.5, -1, 0, 0, 1, 1,
		-0.5, 0.5, -0.5, -1, 0, 0, 0, 1,
		// right
		0.5, -0.5, -0.5, 1, 0, 0, 1, 0,
		0.5, -0.5, 0.5, 1, 0, 0, 0, 0,
		0.5, 0.5, 0.5, 1, 0, 0, 0, 1,
		0.5, 0.5, -0.5, 1, 0, 0, 1, 1,
		// top
		-0.5, 0.5, -0.5, 0, 1, 0, 0, 1,
		0.5, 0.5, -0.5, 0, 1, 0, 1, 1,
		0.5, 0.5, 0.5, 0, 1, 0, 1, 0,
		-0.5, 0.5, 0.5, 0, 1, 0, 0, 0,
		// bottom
		-0.5, -0.5, -0.5, 0, -1, 0, 1, 1,
		0.5, -0.5, -0.5, 0, -1, 0, 0, 1,
		0.5, -0.5, 0.5, 0, -1, 0, 0, 0,
		-0.5, -0.5, 0.5, 0, -1, 0, 1, 0,
	}
	i := []uint32{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
		8, 9, 10, 10, 11, 8,
		12, 13, 14, 14, 15, 12,
		16, 17, 18, 18, 19, 16,
		20, 21, 22, 22, 23, 20,
	}
	return gfx.NewMesh(v, i)
}
