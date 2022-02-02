// Copyright 2020 Nametag, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httperr

import (
	"errors"
	"fmt"
	"io"
)

type errorWithStatus struct {
	statusCode int
	err        error
}

func (e errorWithStatus) Error() string {
	return e.err.Error()
}

func (e errorWithStatus) StatusCode() int {
	return e.statusCode
}

func (e errorWithStatus) Unwrap() error {
	return e.err
}

func (e errorWithStatus) Format(s fmt.State, verb rune) {
	if formatter, ok := e.err.(fmt.Formatter); ok {
		formatter.Format(s, verb)
		return
	}

	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%d %+v", e.StatusCode(), e.err)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// Wrap returns a new error than wraps err embedding statusCode. Calling
// StatusCode() on the return value will return statusCode.
func Wrap(statusCode int, err error) error {
	return errorWithStatus{statusCode: statusCode, err: err}
}

// StatusCoder is an interface for errors that have an integer status
// code attached.
type StatusCoder interface {
	error
	StatusCode() int
}

// StatusCode returns the status code embedded in err, or 0 if no status code is embedded.
func StatusCode(err error) int {
	var sc StatusCoder
	if errors.As(err, &sc) {
		return sc.StatusCode()
	}
	return 0
}
