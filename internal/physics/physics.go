package physics

import "example.com/go-ogl-engine/internal/ecs"

const gravity = -9.81

type World struct {
	Bodies *ecs.Store[ecs.RigidBody]
	Trans  *ecs.Store[ecs.Transform]
}

func New(b *ecs.Store[ecs.RigidBody], t *ecs.Store[ecs.Transform]) *World {
	return &World{Bodies: b, Trans: t}
}

func (p *World) Step(dt float32) {
	p.Bodies.ForEach(func(e ecs.Entity, rb *ecs.RigidBody){
		tr, ok := p.Trans.Get(e)
		if !ok { return }
		if rb.UseGravity {
			rb.Velocity[1] += float32(gravity) * dt
		}
		// Integrate
		tr.Position[0] += rb.Velocity[0] * dt
		tr.Position[1] += rb.Velocity[1] * dt
		tr.Position[2] += rb.Velocity[2] * dt

		// Simple floor collision at y = -1
		if tr.Position[1] < -1.0 {
			tr.Position[1] = -1.0
			rb.Velocity[1] = 0
		}
		p.Trans.Set(e, tr)
	})
}
