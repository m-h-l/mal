package util

type Either[T, U any] struct {
	Left   T
	Right  U
	isLeft bool
}

func NewLeft[T, U any](value T) Either[T, U] {
	return Either[T, U]{Left: value, isLeft: true}
}

func NewRight[T, U any](value U) Either[T, U] {
	return Either[T, U]{Right: value, isLeft: false}
}

func (e Either[T, U]) Get() (T, U) {
	if e.isLeft {
		var zero U
		return e.Left, zero
	} else {
		var zero T
		return zero, e.Right
	}
}

func (e Either[T, U]) GetLeft() (T, bool) {
	if e.isLeft {
		return e.Left, true
	} else {
		var zero T
		return zero, false
	}
}

func (e Either[T, U]) GetRight() (U, bool) {
	if !e.isLeft {
		return e.Right, true
	} else {
		var zero U
		return zero, false
	}
}
func (e Either[T, U]) IsRight() bool {
	return !e.isLeft
}

func (e Either[T, U]) IsLeft() bool {
	return e.isLeft
}
