package xdot

import (
	"errors"
)

type none int

const (
	NONE none = 00
)

func S(err error) func(ErrPipe) {
	return func(pipe ErrPipe) {
		checkX(pipe, err, false)
	}
}

func S1[T any](val T, err error) func(ErrPipe) T {
	return func(pipe ErrPipe) T {
		checkX(pipe, err, false)
		return val
	}
}

func M(err error) func(ErrPipe) {
	return func(pipe ErrPipe) {
		checkX(pipe, err, true)
	}
}

func M1[T any](val T, err error) func(ErrPipe) T {
	return func(pipe ErrPipe) T {
		checkX(pipe, err, true)
		return val
	}
}

func checkX(pipe ErrPipe, err error, quit bool) {
	if err != nil {
		_err := pipe(err)
		if quit {
			panic(_err)
		}
	}
}

type ErrPipe func(err error) error

type Callback func(ErrPipe)

type TryResult[T any] struct {
	_err   error
	Result T
}

func NewTry[T any](result T) *TryResult[T] {
	return &TryResult[T]{
		Result: result,
	}
}

func (tr *TryResult[T]) Return() (T, error) {
	return tr.Result, tr._err
}

func (tr *TryResult[T]) Try(callback Callback) (err error) {
	defer func() {
		r := recover()

		if r == nil {
			return
		}

		_err, ok := r.(error)
		if !ok {
			// not error, re-panic
			panic(r)
		}

		if _err == tr._err {
			err = _err
			return // return this err = _err
		} else {
			// not the same _err, re-panic
			panic(_err)
		}
	}()

	callback(
		// pipe error
		func(err error) error {
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
