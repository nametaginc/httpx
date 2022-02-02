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
)

type errorWithPublic struct {
	public bool
	err    error
}

func (e errorWithPublic) Error() string {
	return e.err.Error()
}

func (e errorWithPublic) IsPublic() bool {
	return e.public
}

func (e errorWithPublic) Unwrap() error {
	return e.err
}

// WrapPublic returns a new error that wraps err, but marked public so that the error content is available in
// http responses.
func WrapPublic(err error) error {
	return errorWithPublic{public: true, err: err}
}

// WrapPrivate returns a new error that wraps err, but marked private so that the error content is not revealed in
// http responses.
func WrapPrivate(err error) error {
	return errorWithPublic{public: false, err: err}
}

// IsPublicer is an interface for errors that implement an IsPublic() method which returns true if the text
// of the error should be reported in an HTTP response by Write()
type IsPublicer interface {
	IsPublic() bool
}

// IsPublic returns true if the error is marked as safe to report via an HTTP response.
func IsPublic(err error) bool {
	var ip IsPublicer
	if errors.As(err, &ip) {
		return ip.IsPublic()
	}
	return false
}

// Publicf is a convenience method that returns a new public error comprising the
// specified message and arguments.
//
// Calling StatusCode() on the return value will return statusCode.
func Publicf(statusCode int, message string, args ...interface{}) error {
	return WrapPublic(Wrapf(statusCode, message, args...))
}

// Wrapf is a convenience method that returns a new error comprising the
// specified message and arguments.
//
// Calling StatusCode() on the return value will return statusCode.
func Wrapf(statusCode int, message string, args ...interface{}) error {
	return errorWithStatus{
		statusCode: statusCode,
		err:        fmt.Errorf(message, args...),
	}
}
