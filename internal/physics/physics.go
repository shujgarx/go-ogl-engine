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
	p.Bodies.ForEach(func(e ecs.Entity, rb *ecs.RigidBody) {
		tr, ok := p.Trans.Get(e)
		if !ok {
			return
		}

		if rb.UseGravity {
			rb.Velocity[1] += float32(gravity) * dt
		}

		if rb.Damping > 0 {
			factor := 1 - rb.Damping*dt
			if factor < 0 {
				factor = 0
			}
			rb.Velocity[0] *= factor
			rb.Velocity[1] *= factor
			rb.Velocity[2] *= factor
		}

		// Integrate simple Euler step
		tr.Position[0] += rb.Velocity[0] * dt
		tr.Position[1] += rb.Velocity[1] * dt
		tr.Position[2] += rb.Velocity[2] * dt

		// Simple floor collision at y = -1 with bounce
		if tr.Position[1] < -1.0 {
			tr.Position[1] = -1.0

			if rb.Velocity[1] < 0 {
				rebound := -rb.Velocity[1] * rb.Bounciness
				if rebound > 0.05 {
					rb.Velocity[1] = rebound
				} else {
					rb.Velocity[1] = 0
				}
			}

			// Apply a little ground friction when we hit the plane
			rb.Velocity[0] *= 0.7
			rb.Velocity[2] *= 0.7
		}

		p.Trans.Set(e, tr)
	})
}
