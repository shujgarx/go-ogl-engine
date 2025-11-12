package engine

import (
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
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

	World      *ecs.World
	Transforms *ecs.Store[ecs.Transform]
	Renderers  *ecs.Store[ecs.MeshRenderer]
	Bodies     *ecs.Store[ecs.RigidBody]

	Phys  *physics.World
	Audio *audio.System

	Meshes   map[string]*gfx.Mesh
	Textures map[string]*gfx.Texture

	Hero ecs.Entity
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
	gl.Enable(gl.CULL_FACE)

	r := &Engine{Window: win, Assets: assets, Meshes: map[string]*gfx.Mesh{}, Textures: map[string]*gfx.Texture{}}
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
	r.Meshes["cube"] = makeCube()
	r.Meshes["floor"] = makeFloor()

	checker, err := gfx.LoadTexture(filepath.Join(assets, "textures", "checker.png"))
	if err != nil {
		return nil, err
	}
	r.Textures["checker"] = checker

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Floating hero cube in the center
	hero := w.NewEntity()
	r.Transforms.Set(hero, ecs.Transform{
		Position: [3]float32{0, 1.25, 0},
		Rotation: [3]float32{0, 0, 0},
		Scale:    [3]float32{1.4, 1.4, 1.4},
	})
	r.Renderers.Set(hero, ecs.MeshRenderer{MeshID: "cube", TextureID: "checker"})
	r.Hero = hero

	// Floor geometry (physics plane handled in physics world)
	floor := w.NewEntity()
	r.Transforms.Set(floor, ecs.Transform{
		Position: [3]float32{0, -1, 0},
		Rotation: [3]float32{0, 0, 0},
		Scale:    [3]float32{22, 1, 22},
	})
	r.Renderers.Set(floor, ecs.MeshRenderer{MeshID: "floor", TextureID: "checker"})

	// Stacks of physics cubes around the hero
	for gx := -1; gx <= 1; gx++ {
		for gz := -1; gz <= 1; gz++ {
			stack := rng.Intn(3) + 3 // 3..5 cubes high
			for gy := 0; gy < stack; gy++ {
				ent := w.NewEntity()
				pos := [3]float32{
					float32(gx)*2.2 + randRange(rng, -0.35, 0.35),
					1.5 + float32(gy)*1.15 + randRange(rng, 0, 0.2),
					float32(gz)*2.2 + randRange(rng, -0.35, 0.35),
				}
				rot := [3]float32{
					randRange(rng, 0, float32(math.Pi*0.25)),
					randRange(rng, 0, float32(math.Pi*2)),
					randRange(rng, 0, float32(math.Pi*0.25)),
				}
				scale := [3]float32{0.75, 0.75, 0.75}
				r.Transforms.Set(ent, ecs.Transform{Position: pos, Rotation: rot, Scale: scale})
				r.Renderers.Set(ent, ecs.MeshRenderer{MeshID: "cube", TextureID: "checker"})

				body := ecs.RigidBody{
					UseGravity: true,
					Mass:       1.0,
					Bounciness: 0.55,
					Damping:    1.2,
				}
				body.Velocity = [3]float32{
					randRange(rng, -1.5, 1.5),
					0,
					randRange(rng, -1.5, 1.5),
				}
				r.Bodies.Set(ent, body)
			}
		}
	}
	return r, nil
}

func (e *Engine) Run() {
	last := time.Now()
	lightPhase := float32(0)
	camPhase := float32(0)
	heroPhase := float32(0)
	fpsTimer := float32(0)
	fpsFrames := 0
	_ = e.Audio.PlaySine(440, 0.2, 0.3)

	for !e.Window.ShouldClose() {
		now := time.Now()
		dt := float32(now.Sub(last).Seconds())
		last = now

		if dt > 0.1 {
			dt = 0.1
		}

		// Update physics
		e.Phys.Step(dt)

		heroPhase += dt
		if e.Hero != 0 {
			if tr, ok := e.Transforms.Get(e.Hero); ok {
				tr.Rotation[1] += dt * 1.4
				tr.Rotation[0] = 0.25 * float32(math.Sin(float64(heroPhase*1.7)))
				tr.Rotation[2] = 0.15 * float32(math.Cos(float64(heroPhase*1.3)))
				e.Transforms.Set(e.Hero, tr)
			}
		}

		lightPhase += dt
		camPhase += dt * 0.35

		lightPos := mgl.Vec3{
			4.0 * float32(math.Cos(float64(lightPhase))),
			3.0 + 1.5*float32(math.Sin(float64(lightPhase*0.5))),
			4.0 * float32(math.Sin(float64(lightPhase))),
		}
		pulse := 0.5 + 0.5*float32(math.Sin(float64(lightPhase*1.5)))
		e.Renderer.LightPosition = lightPos
		e.Renderer.LightColor = mgl.Vec3{
			0.6 + 0.4*pulse,
			0.5 + 0.3*pulse,
			0.8 + 0.2*(1-pulse),
		}

		radius := float32(8)
		camHeight := 2.5 + 1.5*float32(math.Sin(float64(camPhase*0.7)))
		e.Renderer.Cam.Position = mgl.Vec3{
			radius * float32(math.Cos(float64(camPhase))),
			camHeight,
			radius * float32(math.Sin(float64(camPhase))),
		}
		e.Renderer.Cam.Target = mgl.Vec3{0, 0.6, 0}

		// Render
		gl.ClearColor(0.05, 0.07, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(e.Renderer.Program)
		gl.ActiveTexture(gl.TEXTURE0)

		view := e.Renderer.Cam.ViewMatrix()
		proj := e.Renderer.Cam.ProjectionMatrix()
		gl.UniformMatrix4fv(e.Renderer.UniformView, 1, false, &view[0])
		gl.UniformMatrix4fv(e.Renderer.UniformProjection, 1, false, &proj[0])

		lp := e.Renderer.LightPosition
		lc := e.Renderer.LightColor
		gl.Uniform3f(e.Renderer.UniformLightPos, lp[0], lp[1], lp[2])
		gl.Uniform3f(e.Renderer.UniformLightColor, lc[0], lc[1], lc[2])
		camPos := e.Renderer.Cam.Position
		gl.Uniform3f(e.Renderer.UniformViewPos, camPos[0], camPos[1], camPos[2])

		e.Transforms.ForEach(func(ent ecs.Entity, t *ecs.Transform) {
			if mr, ok := e.Renderers.Get(ent); ok {
				mesh := e.Meshes[mr.MeshID]
				if mesh == nil {
					return
				}
				if tex := e.Textures[mr.TextureID]; tex != nil {
					gl.BindTexture(gl.TEXTURE_2D, tex.ID)
				} else {
					gl.BindTexture(gl.TEXTURE_2D, 0)
				}
				model := gfx.ModelMatrix(t.Position, t.Rotation, t.Scale)
				normal := gfx.NormalMatrix(model)
				gl.UniformMatrix4fv(e.Renderer.UniformModel, 1, false, &model[0])
				gl.UniformMatrix3fv(e.Renderer.UniformNormal, 1, false, &normal[0])
				mesh.Draw()
			}
		})

		e.Window.SwapBuffers()
		glfw.PollEvents()

		fpsTimer += dt
		fpsFrames++
		if fpsTimer >= 0.5 {
			fps := float64(fpsFrames) / float64(fpsTimer)
			e.Window.SetTitle(fmt.Sprintf("Go OGL Engine Tech Demo | %.0f FPS", fps))
			fpsTimer = 0
			fpsFrames = 0
		}
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

func makeFloor() *gfx.Mesh {
	tile := float32(10)
	v := []float32{
		// positions         // normals        // uvs
		-0.5, 0, 0.5, 0, 1, 0, 0, tile,
		0.5, 0, 0.5, 0, 1, 0, tile, tile,
		0.5, 0, -0.5, 0, 1, 0, tile, 0,
		-0.5, 0, -0.5, 0, 1, 0, 0, 0,
	}
	i := []uint32{0, 1, 2, 2, 3, 0}
	return gfx.NewMesh(v, i)
}

func randRange(r *rand.Rand, min, max float32) float32 {
	return min + r.Float32()*(max-min)
}
