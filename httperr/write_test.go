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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

type differentResponseWriter struct {
}

func (differentResponseWriter) Error() string {
	return "foo"
}

func (differentResponseWriter) WriteResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(999)
	fmt.Fprintln(w, "Hello, world!")
}

func TestWrite(t *testing.T) {
	// private
	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		err := Wrap(418, fmt.Errorf("cannot frob the grob"))
		Write(w, r, err)
		assert.Check(t, is.Equal(418, w.Code))
		assert.Check(t, is.Equal("I'm a teapot\n", string(w.Body.Bytes())))
	}

	// public
	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		err := WrapPublic(Wrap(418, fmt.Errorf("cannot frob the grob")))
		Write(w, r, err)
		assert.Check(t, is.Equal(418, w.Code))
		assert.Check(t, is.Equal("I'm a teapot\n", string(w.Body.Bytes())))
	}

	// custom
	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/", nil)
		err := errors.Wrap(differentResponseWriter{}, "foo")
		Write(w, r, err)
		assert.Check(t, is.Equal(999, w.Code))
		assert.Check(t, is.Equal("Hello, world!\n", string(w.Body.Bytes())))
	}
}
