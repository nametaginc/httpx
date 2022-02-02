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
	"context"
	"net/http"
)

type onErrorIndexType int

const onErrorIndex onErrorIndexType = iota

// OnError returns a new http.Request that holds a reference to a
// function that will report an error when returned from a request.
func OnError(r *http.Request, f func(err error)) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), onErrorIndex, f))
}

// ReportError reports the error to the function given in
// OnError.
func ReportError(w http.ResponseWriter, r *http.Request, err error) {
	err = TranslateError(err)

	if v := r.Context().Value(onErrorIndex); v != nil {
		v.(func(error))(err)
	} else {
		Write(w, r, err)
	}
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// handler object that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		ReportError(w, r, err)
	}
}
