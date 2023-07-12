package result

import "fmt"

type Result[T any] struct {
	value *T
	error error
}

func Value[T any](value T) Result[T] {
	return Result[T]{value: &value}
}

func Error[T any](err error) Result[T] {
	return Result[T]{error: err}
}

func ErrorWrap[T any](err error, msg string) Result[T] {
	return Result[T]{
		error: fmt.Errorf("%s: %w", msg, err),
	}
}

func (r *Result[T]) Value() (T, error) {
	if r.error != nil {
		var t T
		return t, r.error
	}
	return *r.value, nil
}
