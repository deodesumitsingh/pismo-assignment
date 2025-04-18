package v1

type Handler[T any] struct {
	service T
}

func NewHandler[T any](service T) *Handler[T] {
	return &Handler[T]{
		service: service,
	}
}
