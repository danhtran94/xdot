package x

import "errors"

type ErrHandler func(error)

// TryUnsafe is a function that returns a defer function and a pipe function. The defer function is used to handle the error,
// use it when you want to handle the error in the defer block instead of the nested block of try's call block.
//
//	try, pipe := TryPipe()
//	defer try(func(err error) {
//		fmt.Println(err)
//	})
func TryUnsafe() (deferTry func(...ErrHandler), pipe ErrPipe) {
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
