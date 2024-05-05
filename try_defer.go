package xdot

import "errors"

type ErrHandler func(error)

// TryPipe Example:
//
//	try, pipe := TryPipe()
//	defer try(func(err error) {
//		fmt.Println(err)
//	})
func TryPipe() (deferTry func(...ErrHandler), pipe ErrPipe) {
	var _err error

	pipe = func(err error) error {
		if err != nil {
			if _err == nil {
				_err = err
			} else {
				_err = errors.Join(_err, err)
			}
		}

		return _err
	}

	deferTry = func(handlers ...ErrHandler) {
		r := recover()
		if r == nil {
			return
		}

		err, ok := r.(error)
		if !ok {
			panic(r) // not error, re-panic
		}

		if err == _err {
			for _, h := range handlers {
				h(_err)
			}
			return // return this err = _err
		} else {
			panic(err) // not the same _err, re-panic
		}
	}

	return deferTry, pipe
}

func Must0(err error) func(ErrPipe) {
	return func(pipe ErrPipe) {
		checkX(pipe, err, true)
	}
}

func Must[T any](val T, err error) func(ErrPipe) T {
	return func(pipe ErrPipe) T {
		checkX(pipe, err, true)
		return val
	}
}

func Should0(err error) func(ErrPipe) {
	return func(pipe ErrPipe) {
		checkX(pipe, err, false)
	}
}

func Should[T any](val T, err error) func(ErrPipe) T {
	return func(pipe ErrPipe) T {
		checkX(pipe, err, false)
		return val
	}
}
