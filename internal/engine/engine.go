package engine

import (
	"fmt"
	"runtime"
	"time"

	"example.com/go-ogl-engine/internal/audio"
	"example.com/go-ogl-engine/internal/ecs"
	"example.com/go-ogl-engine/internal/gfx"
	"example.com/go-ogl-engine/internal/physics"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	Window   *glfw.Window
	Renderer *gfx.Renderer
	Assets   string

	Width  int
	Height int

	ClearColor [4]float32

	World      *ecs.World
	Transforms *ecs.Store[ecs.Transform]
	Renderers  *ecs.Store[ecs.MeshRenderer]
	Bodies     *ecs.Store[ecs.RigidBody]

	Phys  *physics.World
	Audio *audio.System

	Meshes   map[string]*gfx.Mesh
	Textures map[string]*gfx.Texture

	project Project
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
	gl.Enable(gl.CULL_FACE) // ФИКСИТЬ, с gl.Disable(gl.CULL_FACE) Все корректно отображается, однако если включить то полигоны у пирамиды пропадают, и надо проверить остальные фигуры

	fbw, fbh := win.GetFramebufferSize()
	r := &Engine{
		Window:     win,
		Assets:     assets,
		Width:      fbw,
		Height:     fbh,
		ClearColor: [4]float32{0.05, 0.07, 0.1, 1.0},
		Meshes:     map[string]*gfx.Mesh{},
		Textures:   map[string]*gfx.Texture{},
		project:    NullProject{},
	}

	r.Renderer, err = gfx.NewRenderer(assets, fbw, fbh)
	if err != nil {
		return nil, err
	}

	win.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		if width <= 0 || height <= 0 {
			return
		}
		r.Width = width
		r.Height = height
		r.Renderer.Resize(width, height)
	})

	w := ecs.NewWorld()
	r.World = w
	r.Transforms = ecs.NewStore[ecs.Transform]()
	r.Renderers = ecs.NewStore[ecs.MeshRenderer]()
	r.Bodies = ecs.NewStore[ecs.RigidBody]()

	r.Phys = physics.New(r.Bodies, r.Transforms)

	r.Audio = audio.NewSystem()
	_ = r.Audio.Init(48000)

	return r, nil
}

func (e *Engine) UseProject(p Project) {
	e.project = p
}

func (e *Engine) Run() {
	if e.project == nil {
		return
	}
	if err := e.project.Load(e); err != nil {
		fmt.Println("project load error:", err)
		return
	}
	defer e.project.Unload(e)

	last := time.Now()
	fpsTimer := float32(0)
	fpsFrames := 0

	for !e.Window.ShouldClose() {
		now := time.Now()
		dt := float32(now.Sub(last).Seconds())
		last = now

		if dt > 0.1 {
			dt = 0.1
		}

		e.Phys.Step(dt)
		e.World.Update(dt)

		e.project.Update(e, dt)

		e.Renderer.Cam.Aspect = float32(e.Width) / float32(e.Height)

		e.Renderer.BeginShadowPass()
		e.forEachRenderable(func(mesh *gfx.Mesh, _ *gfx.Texture, model mgl.Mat4, _ mgl.Mat3) {
			e.Renderer.SubmitShadow(model)
			mesh.Draw()
		})
		e.Renderer.EndShadowPass()

		e.Renderer.BeginMainPass(e.Width, e.Height, e.ClearColor)
		e.forEachRenderable(func(mesh *gfx.Mesh, tex *gfx.Texture, model mgl.Mat4, normal mgl.Mat3) {
			if tex != nil {
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, tex.ID)
			} else {
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, 0)
			}
			e.Renderer.SubmitMesh(model, normal)
			mesh.Draw()
		})
		e.Renderer.EndMainPass()

		e.project.Render(e)

		e.Window.SwapBuffers()
		glfw.PollEvents()

		fpsTimer += dt
		fpsFrames++
		if fpsTimer >= 0.5 {
			fps := float64(fpsFrames) / float64(fpsTimer)
			e.Window.SetTitle(fmt.Sprintf("Go OGL Engine | %.0f FPS", fps))
			fpsTimer = 0
			fpsFrames = 0
		}
	}
	e.Shutdown()
}

func (e *Engine) forEachRenderable(fn func(*gfx.Mesh, *gfx.Texture, mgl.Mat4, mgl.Mat3)) {
	e.Transforms.ForEach(func(ent ecs.Entity, t *ecs.Transform) {
		if mr, ok := e.Renderers.Get(ent); ok {
			mesh := e.Meshes[mr.MeshID]
			if mesh == nil {
				return
			}
			var tex *gfx.Texture
			if mr.TextureID != "" {
				tex = e.Textures[mr.TextureID]
			}
			model := gfx.ModelMatrix(t.Position, t.Rotation, t.Scale)
			normal := gfx.NormalMatrix(model)
			fn(mesh, tex, model, normal)
		}
	})
}

func (e *Engine) Shutdown() {
	glfw.Terminate()
}
