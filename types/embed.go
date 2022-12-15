package types

type Embed[T any] interface {
	Get() []T
	Add(val T)
}

func NewEmbed[T any]() EmbedTemplate[T] {
	return EmbedTemplate[T]{
		Slices: make([]T, 0),
	}
}

type EmbedTemplate[T any] struct {
	Slices []T
}

func (em *EmbedTemplate[T]) Get() []T {
	return em.Slices
}

func (em *EmbedTemplate[T]) Add(val T) {
	em.Slices = append(em.Slices, val)
}
