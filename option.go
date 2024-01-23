package common

type Option[T any] struct {
	internal T
	set      bool
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func Some[T any](t T) Option[T] {
	return Option[T]{
		internal: t,
		set:      true,
	}
}

func (opt *Option[T]) Ok() bool {
	return opt.set
}

func (opt *Option[T]) Unwrap() T {
	if opt.set {
		return opt.internal
	}

	panic("unwrap called on none")
}
