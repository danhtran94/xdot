package x

import (
	"errors"
	"fmt"
)

type none int
type ErrPipe func(err error) error

func (p ErrPipe) With(err error) ErrPipe {
	return func(orig error) error {
		return p(wrapErrs(orig, err))
	}
}

func (p ErrPipe) When(trigger bool) ErrPipe {
	return func(err error) error {
		if trigger {
			return p(err)
		}

		return p(nil)
	}
}

func (p ErrPipe) Err() error {
	return p(nil)
}

type Callback func(ErrPipe)
type TryResult[T any] struct {
	_err   error
	Result T
}

const (
	NONE none = 00
)

func checkX(pipe ErrPipe, err error, quit bool) {
	if err != nil {
		_err := pipe(err)
		if quit && _err != nil {
			panic(_err)
		}
	}
}

func Try[T any](result T) *TryResult[T] {
	return &TryResult[T]{
		Result: result,
	}
}

func (tr *TryResult[T]) Return() (T, error) {
	return tr.Result, tr._err
}

func (tr *TryResult[T]) Call(callback Callback) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		_err, ok := r.(error)
		if !ok {
			panic(r) // not error, re-panic
		}

		if _err == tr._err {
			err = _err
			return // return this err = _err
		} else {
			panic(_err) // not the same _err, re-panic
		}
	}()

	callback(
		func(err error) error { // pipe error
			if err != nil {
				if tr._err == nil {
					tr._err = err
				} else {
					tr._err = errors.Join(tr._err, err)
				}
			}

			return tr._err
		})

	return tr._err
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

func wrapErrs(orig error, err error) error {
	if err == nil {
		return orig
	}

	if orig != nil {
		return fmt.Errorf("%w; %w", err, orig)
	}

	return nil
}
