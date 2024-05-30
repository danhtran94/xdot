package xdot

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
)

func Wrap(orig error, err error) error {
	if err == nil {
		return orig
	}

	if orig != nil {
		return fmt.Errorf("%w; %w", err, orig)
	}

	return nil
}

func LogErr(err error, errs ...error) error {
	if err == nil {
		return nil
	}

	log.Printf("[ERR] %v\n=== Stack TraceBack ===\n%s=== Stack EndTrace ===\n", err, debug.Stack())
	return Wrap(err, errors.Join(errs...))
}

func HandleLogErr(err error) ErrHandler {
	return func(orig error) {
		LogErr(orig, err)
	}
}
