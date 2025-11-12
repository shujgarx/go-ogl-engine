package ecs

import (
	"container/list"
	"sync"

	"github.com/ebitengine/oto/v3"
)

type Entity int

// System интерфейс, который должны реализовывать все системы
type System interface {
	Update(w *World, dt float32)
}

// AudioSystem — система звука
type AudioSystem struct {
	ctx *oto.Context
	mu  sync.Mutex
	// Можно добавить сюда, например, хранилище AudioSource, или ссылки на сущности и т.д.
}

func NewAudioSystem(ctx *oto.Context) *AudioSystem {
	return &AudioSystem{ctx: ctx}
}

// Реализация метода Update для AudioSystem
func (s *AudioSystem) Update(w *World, dt float32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Здесь логика обновления звука, например проверка компонентов AudioSource и проигрывание/остановка
	// Пока просто пример:
	// fmt.Println("AudioSystem Update dt =", dt)
}

type World struct {
	next    Entity
	Alive   map[Entity]bool
	Systems []System
}

func NewWorld() *World {
	return &World{
		Alive:   map[Entity]bool{},
		Systems: []System{},
	}
}

func (w *World) NewEntity() Entity {
	w.next++
	w.Alive[w.next] = true
	return w.next
}

func (w *World) AddSystem(s System) {
	w.Systems = append(w.Systems, s)
}

func (w *World) Update(dt float32) {
	for _, s := range w.Systems {
		s.Update(w, dt)
	}
}

// --- Components ---
type Transform struct {
	Position [3]float32
	Rotation [3]float32
	Scale    [3]float32
}

type MeshRenderer struct {
	MeshID    string
	TextureID string
}

type RigidBody struct {
	Velocity   [3]float32
	Mass       float32
	UseGravity bool
}

type AudioSource struct {
	Playing bool
	FreqHz  float64
	Gain    float64
}

// Component stores
type Store[T any] struct {
	Data  map[Entity]T
	Order *list.List
}

func NewStore[T any]() *Store[T] {
	return &Store[T]{Data: map[Entity]T{}, Order: list.New()}
}

func (s *Store[T]) Set(e Entity, v T) {
	s.Data[e] = v
	s.Order.PushBack(e)
}

func (s *Store[T]) Get(e Entity) (T, bool) {
	v, ok := s.Data[e]
	return v, ok
}

func (s *Store[T]) Has(e Entity) bool {
	_, ok := s.Data[e]
	return ok
}

func (s *Store[T]) ForEach(fn func(e Entity, v *T)) {
	for el := s.Order.Front(); el != nil; el = el.Next() {
		e := el.Value.(Entity)
		v := s.Data[e]
		fn(e, &v)
		s.Data[e] = v
	}
}
