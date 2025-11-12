package daylight

import (
	"math"
	"path/filepath"

	"example.com/go-ogl-engine/internal/ecs"
	"example.com/go-ogl-engine/internal/engine"
	"example.com/go-ogl-engine/internal/gfx"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Demo struct {
	time    float32
	target  mgl.Vec3
	sphere  ecs.Entity
	pyramid ecs.Entity
	cone    ecs.Entity
}

func New() *Demo {
	return &Demo{target: mgl.Vec3{0, 0, 0}}
}

func (d *Demo) Load(e *engine.Engine) error {
	if _, ok := e.Meshes["cube"]; !ok {
		e.Meshes["cube"] = gfx.NewCubeMesh()
	}
	if _, ok := e.Meshes["floor"]; !ok {
		e.Meshes["floor"] = gfx.NewPlaneMesh(12)
	}
	if _, ok := e.Meshes["sphere"]; !ok {
		e.Meshes["sphere"] = gfx.NewSphereMesh(32, 32)
	}
	if _, ok := e.Meshes["pyramid"]; !ok {
		e.Meshes["pyramid"] = gfx.NewPyramidMesh()
	}
	if _, ok := e.Meshes["cone"]; !ok {
		e.Meshes["cone"] = gfx.NewConeMesh(32)
	}

	if _, ok := e.Textures["checker"]; !ok {
		tex, err := gfx.LoadTexture(filepath.Join(e.Assets, "textures", "checker.png"))
		if err != nil {
			return err
		}
		e.Textures["checker"] = tex
	}

	e.ClearColor = [4]float32{0.62, 0.78, 0.95, 1.0}

	floor := e.World.NewEntity()
	e.Transforms.Set(floor, ecs.Transform{
		Position: [3]float32{0, -1, 0},
		Rotation: [3]float32{0, 0, 0},
		Scale:    [3]float32{32, 1, 32},
	})
	e.Renderers.Set(floor, ecs.MeshRenderer{MeshID: "floor", TextureID: "checker"})

	cube := e.World.NewEntity()
	e.Transforms.Set(cube, ecs.Transform{
		Position: [3]float32{-3, -0.2, -2},
		Rotation: [3]float32{0, 0.2, 0},
		Scale:    [3]float32{1.5, 1.5, 1.5},
	})
	e.Renderers.Set(cube, ecs.MeshRenderer{MeshID: "cube", TextureID: "checker"})

	d.sphere = e.World.NewEntity()
	e.Transforms.Set(d.sphere, ecs.Transform{
		Position: [3]float32{-4, 0, 3},
		Rotation: [3]float32{0, 0, 0},
		Scale:    [3]float32{2, 2, 2},
	})
	e.Renderers.Set(d.sphere, ecs.MeshRenderer{MeshID: "sphere", TextureID: "checker"})

	d.pyramid = e.World.NewEntity()
	e.Transforms.Set(d.pyramid, ecs.Transform{
		Position: [3]float32{0, 0, 0},
		Rotation: [3]float32{0, 0, 0},
		Scale:    [3]float32{2.2, 2.2, 2.2},
	})
	e.Renderers.Set(d.pyramid, ecs.MeshRenderer{MeshID: "pyramid", TextureID: "checker"})

	d.cone = e.World.NewEntity()
	e.Transforms.Set(d.cone, ecs.Transform{
		Position: [3]float32{4, 0, -1.5},
		Rotation: [3]float32{0, 0, 0},
		Scale:    [3]float32{2, 2, 2},
	})
	e.Renderers.Set(d.cone, ecs.MeshRenderer{MeshID: "cone", TextureID: "checker"})

	e.Renderer.Cam.Position = mgl.Vec3{10, 6, 10}
	e.Renderer.Cam.Target = mgl.Vec3{0, 0.5, 0}

	e.Renderer.LightColor = mgl.Vec3{1.0, 0.97, 0.9}
	return nil
}

func (d *Demo) Update(e *engine.Engine, dt float32) {
	d.time += dt

	sunAngle := d.time * 0.25
	lightDir := mgl.Vec3{
		float32(math.Cos(float64(sunAngle))) * 0.6,
		-1.0,
		float32(math.Sin(float64(sunAngle))) * 0.6,
	}.Normalize()
	lightPos := d.target.Sub(lightDir.Mul(20))
	e.Renderer.LightDirection = lightDir
	e.Renderer.LightPosition = lightPos

	lightView := mgl.LookAtV(lightPos, d.target, mgl.Vec3{0, 1, 0})
	lightProj := mgl.Ortho(-25, 25, -25, 25, 1, 60)
	e.Renderer.LightSpaceMatrix = lightProj.Mul4(lightView)

	radius := float32(16)
	camAngle := d.time * 0.15
	e.Renderer.Cam.Position = mgl.Vec3{
		radius * float32(math.Cos(float64(camAngle))),
		7.0,
		radius * float32(math.Sin(float64(camAngle))),
	}
	e.Renderer.Cam.Target = mgl.Vec3{0, 0.8, 0}

	if tr, ok := e.Transforms.Get(d.sphere); ok {
		tr.Position[1] = float32(math.Sin(float64(d.time*1.2))) * 0.5
		tr.Rotation[1] += dt * 0.5
		e.Transforms.Set(d.sphere, tr)
	}
	if tr, ok := e.Transforms.Get(d.pyramid); ok {
		tr.Rotation[1] += dt * 0.8
		e.Transforms.Set(d.pyramid, tr)
	}
	if tr, ok := e.Transforms.Get(d.cone); ok {
		tr.Rotation[0] += dt * 0.6
		tr.Rotation[2] += dt * 0.4
		e.Transforms.Set(d.cone, tr)
	}
}

func (d *Demo) Render(*engine.Engine) {}

func (d *Demo) Unload(*engine.Engine) {}
