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
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

type differentStatusCoder struct {
}

func (differentStatusCoder) Error() string {
	return "foo"
}

func (differentStatusCoder) StatusCode() int {
	return 123
}

func removePath(s string) string {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	s = strings.ReplaceAll(s, curDir, "<PACKAGE_DIR>")
	s = strings.ReplaceAll(s, runtime.GOROOT(), "<GOROOT>")
	s = strings.ReplaceAll(s, runtime.GOARCH, "<GOARCH>")
	s = regexp.MustCompile(`.(go|s):\d+`).ReplaceAllString(s, ".${1}:<LINE>")
	return s
}

func TestWrap(t *testing.T) {
	e1 := fmt.Errorf("cannot frob the grob")
	e2 := Wrap(419, e1)
	e3 := errors.Wrap(e2, "while trying to frob the grov")
	e4 := Wrap(999, e3)

	assert.Check(t, is.Equal(0, StatusCode(e1)))
	assert.Check(t, is.Equal(419, StatusCode(e2)))
	assert.Check(t, is.Equal(419, StatusCode(e3)))
	assert.Check(t, is.Equal(999, StatusCode(e4)))

	e5 := differentStatusCoder{}
	assert.Check(t, is.Equal(123, StatusCode(e5)))

	// e3 shows the stack with +v
	assert.Check(t, is.Contains(removePath(fmt.Sprintf("%+v", e3)), "/httpx/httperr"))

	// e4 wraps an error with a stack, and the stack is still visible
	assert.Check(t, is.Contains(removePath(fmt.Sprintf("%+v", e4)), "/httpx/httperr"))
}
