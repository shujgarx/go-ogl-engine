package engine

type Project interface {
	Load(*Engine) error
	Update(*Engine, float32)
	Render(*Engine)
	Unload(*Engine)
}

// NullProject is a no-op implementation used as default placeholders.
type NullProject struct{}

func (NullProject) Load(*Engine) error      { return nil }
func (NullProject) Update(*Engine, float32) {}
func (NullProject) Render(*Engine)          {}
func (NullProject) Unload(*Engine)          {}
