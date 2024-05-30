package xdot

import (
	"errors"
)

type none int
type ErrPipe func(err error) error

func (p ErrPipe) With(err error) ErrPipe {
	return func(orig error) error {
		return p(Wrap(orig, err))
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

func I[T any](val T, err error) func() (T, error) {
	return func() (T, error) {
		return val, err
	}
}

func O(err error) func() (error, error) {
	return func() (error, error) {
		return err, err
	}
}

func M[T any](fn func() (T, error), pipe ErrPipe) T {
	val, err := fn()
	checkX(pipe, err, true)
	return val
}

func S[T any](fn func() (T, error), pipe ErrPipe) T {
	val, err := fn()
	checkX(pipe, err, false)
	return val
}

func checkX(pipe ErrPipe, err error, quit bool) {
	if err != nil {
		_err := pipe(err)
		if quit {
			panic(_err)
		}
	}
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
