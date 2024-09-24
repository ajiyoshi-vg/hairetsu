package result

type Result[T any] struct {
	ok  T
	err error
}

func (r *Result[T]) OK() T {
	return r.ok
}

func (r *Result[T]) Error() error {
	return r.err
}

func (r *Result[T]) Result() (T, error) {
	return r.ok, r.err
}

func New[T any](x T, err error) Result[T] {
	if err != nil {
		return NG[T](err)
	}
	return OK(x)
}

func OK[T any](ok T) Result[T] {
	return Result[T]{ok: ok}
}

func NG[T any](err error) Result[T] {
	return Result[T]{err: err}
}
