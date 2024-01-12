# xdot Go Library

**The error handling library that missing in Go.**

[![Go Reference](https://pkg.go.dev/badge/github.com/danhtran94/xdot.svg)](https://pkg.go.dev/github.com/danhtran94/xdot)
[![Go Report Card](https://goreportcard.com/badge/github.com/danhtran94/xdot)](https://goreportcard.com/report/github.com/danhtran94/xdot)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This library provides a set of functions and types to handle errors in Go in a more generic functional way. It allows you to pipe errors through functions and handle them in a centralized manner.

**Super tiny library, with only 1 file, < 100 LOCs and 0 dependencies.** \
You can read through the whole library in **5 minutes**.

**Less code, less bugs, save time, save money.**

## Installation

To install this library, use the `go get` command:

```bash
go get github.com/danhtran94/xdot
```

## Example

Here is a basic example of how to use this library:
It can be used in more complex cases your contribution examples are welcome.

```go
package main

import (
	"io"
	"net/http"

	. "github.com/danhtran94/xdot"
)

// Comparision With and Without xdot: Get Github repo's branch main information

// Total LOCs: 11 lines
func WithXDot(repo string) ([]byte, error) {
	tr := NewTry([]byte{})
	tr.Try(func(pipe ErrPipe) {
		branchMainURL := "https://api.github.com/repos/" + repo + "/branches/main" // Declare (1 line)
		req := M(I(http.NewRequest("GET", branchMainURL, nil)), pipe)              // Logic (4 lines)
		body := M(I(http.DefaultClient.Do(req)), pipe).Body
		defer body.Close()
		tr.Result = M(I(io.ReadAll(body)), pipe)
	})
	return tr.Return() // Return (1 line)
}

// Total LOCs: 18 lines (> 63.6% ~ 7 lines more than WithXDot version)
func WithoutXDot(repo string) ([]byte, error) {
	branchMainURL := "https://api.github.com/repos/" + repo + "/branches/main" // Declare (1 line)
	req, err := http.NewRequest("GET", branchMainURL, nil)                     // Logic (14 lines)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer body.Close()
	result, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return result, nil // Return (1 line)
}

```


## Usage

Here is a brief overview of the main components of the library:

### Types

- `TryResult`: A generic struct that holds a result and an error. It has methods to handle errors in a functional way.
- `Callback`: A function type that takes an `ErrPipe` function.
- `ErrPipe`: A function type that takes an error and returns an error. It's used to pipe errors through functions.
- `NONE`: Used to indicate that a function does not return a value.

### Functions

- `NewTry` A function that creates a new `TryResult` with a given result.
- `M` means **"must"** execute function successfully if not flow will be stopped and "error" will be returned.
- `S` means **"should"** execute function successfully if not flow will be continued and "error" will saved in TryResult.
- `O` and `I` are used to indicate whether a function returns a value or not.
- `Try` A method of `TryResult` that takes a `Callback` function and handles errors in a functional way.
- `Return`: A method of `TryResult` that returns the result and the error.

## Contributing

Contributions are welcome. Please submit a pull request or create an issue to discuss the changes you want to make.

## License

This library is licensed under the MIT License.
