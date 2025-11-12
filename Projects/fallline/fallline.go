package fallline

import (
	"math"
	"math/rand"
	"path/filepath"
	"time"

	"example.com/go-ogl-engine/internal/ecs"
	"example.com/go-ogl-engine/internal/engine"
	"example.com/go-ogl-engine/internal/gfx"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type drop struct {
	entity ecs.Entity
	spin   [3]float32
}

type Demo struct {
	time         float32
	rng          *rand.Rand
	drops        []*drop
	nextSpawnZ   float32
	spawnSpacing float32
	cameraZ      float32
	lightDir     mgl.Vec3
	meshes       []string
	textures     []string
}

func New() *Demo {
	return &Demo{
		nextSpawnZ:   18,
		spawnSpacing: 18,
		cameraZ:      -36,
		lightDir:     mgl.Vec3{-0.35, -1.0, -0.35}.Normalize(),
	}
}

func (d *Demo) Load(e *engine.Engine) error {
	d.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	d.meshes = []string{"cube", "sphere", "pyramid", "cone"}

	if _, ok := e.Meshes["cube"]; !ok {
		e.Meshes["cube"] = gfx.NewCubeMesh()
	}
	if _, ok := e.Meshes["sphere"]; !ok {
		e.Meshes["sphere"] = gfx.NewSphereMesh(24, 24)
	}
	if _, ok := e.Meshes["pyramid"]; !ok {
		e.Meshes["pyramid"] = gfx.NewPyramidMesh()
	}
	if _, ok := e.Meshes["cone"]; !ok {
		e.Meshes["cone"] = gfx.NewConeMesh(32)
	}
	if _, ok := e.Meshes["floor"]; !ok {
		e.Meshes["floor"] = gfx.NewPlaneMesh(64)
	}

	textureSources := map[string]string{
		"checker":    filepath.Join(e.Assets, "textures", "checker.png"),
		"clay_color": filepath.Join(e.Assets, "materials", "clay_smooth_whitebrown", "Ground026_1K-PNG_Color.png"),
		"clay_core":  filepath.Join(e.Assets, "materials", "clay_smooth_whitebrown", "Ground026.png"),
		"concrete":   filepath.Join(e.Assets, "materials", "concrete_plaster_white", "Concrete034_1K-PNG_Color.png"),
	}
	for id, src := range textureSources {
		if _, ok := e.Textures[id]; ok {
			continue
		}
		tex, err := gfx.LoadTexture(src)
		if err != nil {
			return err
		}
		e.Textures[id] = tex
	}

	d.textures = []string{"checker", "checker", "clay_color", "clay_color", "clay_core", "concrete", "concrete", "clay_color", "checker", "concrete"}

	e.ClearColor = [4]float32{0.03, 0.04, 0.06, 1.0}

	const segmentLength = float32(80)
	const segmentCount = 6
	for i := 0; i < segmentCount; i++ {
		seg := e.World.NewEntity()
		centerZ := (float32(i) - 1) * segmentLength
		e.Transforms.Set(seg, ecs.Transform{
			Position: [3]float32{0, -1, centerZ},
			Rotation: [3]float32{0, 0, 0},
			Scale:    [3]float32{16, 1, segmentLength},
		})
		texID := "checker"
		if i%2 == 1 {
			texID = "clay_color"
		}
		e.Renderers.Set(seg, ecs.MeshRenderer{MeshID: "floor", TextureID: texID})
	}

	e.Renderer.LightColor = mgl.Vec3{1.0, 0.95, 0.85}

	const dropCount = 10
	for i := 0; i < dropCount; i++ {
		d.drops = append(d.drops, d.spawnDrop(e))
	}

	e.Renderer.Cam.Position = mgl.Vec3{0, 3, d.cameraZ}
	e.Renderer.Cam.Target = mgl.Vec3{0, 2.5, d.cameraZ + 5}

	return nil
}

func (d *Demo) spawnDrop(e *engine.Engine) *drop {
	ent := e.World.NewEntity()
	e.Bodies.Set(ent, ecs.RigidBody{
		UseGravity: true,
		Mass:       1,
		Bounciness: 0.35,
		Damping:    0.08,
	})
	dr := &drop{entity: ent}
	d.resetDrop(e, dr)
	return dr
}

func (d *Demo) resetDrop(e *engine.Engine, dr *drop) {
	meshID := d.meshes[d.rng.Intn(len(d.meshes))]
	texID := d.textures[d.rng.Intn(len(d.textures))]

	scale := 0.8 + d.rng.Float32()*0.7
	x := (d.rng.Float32()*9 - 4.5)
	y := 5 + d.rng.Float32()*3
	z := d.nextSpawnZ + d.rng.Float32()*6
	d.nextSpawnZ += d.spawnSpacing

	e.Renderers.Set(dr.entity, ecs.MeshRenderer{MeshID: meshID, TextureID: texID})
	e.Transforms.Set(dr.entity, ecs.Transform{
		Position: [3]float32{x, y, z},
		Rotation: [3]float32{0, 0, 0},
		Scale:    [3]float32{scale, scale, scale},
	})

	if rb, ok := e.Bodies.Get(dr.entity); ok {
		rb.Velocity = [3]float32{0, 0, 0}
		e.Bodies.Set(dr.entity, rb)
	}

	dr.spin = [3]float32{
		(d.rng.Float32() - 0.5) * 1.4,
		(d.rng.Float32() - 0.5) * 1.4,
		(d.rng.Float32() - 0.5) * 1.0,
	}
}

func (d *Demo) Update(e *engine.Engine, dt float32) {
	d.time += dt
	speed := float32(10)
	d.cameraZ += speed * dt

	camHeight := float32(3.2 + 0.3*math.Sin(float64(d.time*0.7)))
	camPos := mgl.Vec3{0, camHeight, d.cameraZ}
	lookAhead := float32(8)
	target := mgl.Vec3{0, camHeight - 0.6, d.cameraZ + lookAhead}

	e.Renderer.Cam.Position = camPos
	e.Renderer.Cam.Target = target

	focus := mgl.Vec3{0, -0.5, d.cameraZ + lookAhead}
	lightPos := focus.Sub(d.lightDir.Normalize().Mul(60))
	e.Renderer.LightDirection = d.lightDir
	e.Renderer.LightPosition = lightPos
	lightProj := mgl.Ortho(-35, 35, -35, 35, 1, 130)
	lightView := mgl.LookAtV(lightPos, focus, mgl.Vec3{0, 1, 0})
	e.Renderer.LightSpaceMatrix = lightProj.Mul4(lightView)

	for _, dr := range d.drops {
		tr, ok := e.Transforms.Get(dr.entity)
		if !ok {
			continue
		}
		tr.Rotation[0] += dr.spin[0] * dt
		tr.Rotation[1] += dr.spin[1] * dt
		tr.Rotation[2] += dr.spin[2] * dt
		e.Transforms.Set(dr.entity, tr)

		if tr.Position[2] < d.cameraZ-18 || tr.Position[1] < -1.2 {
			d.resetDrop(e, dr)
		}
	}
}

func (d *Demo) Render(*engine.Engine) {}

func (d *Demo) Unload(*engine.Engine) {}
