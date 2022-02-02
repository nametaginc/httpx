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

package httpx

import (
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

type headerTest struct {
	in   string
	p, v string
}

func TestSplitHeaderOnSpace(t *testing.T) {
	tests := []headerTest{
		{in: "biscuit value", p: "biscuit", v: "value"},
		{in: "biscuit    value", p: "biscuit", v: "value"},
		{in: "biscuit   biscuit two", p: "biscuit", v: "biscuit two"},
		{in: "bearer  value with spaces", p: "bearer", v: "value with spaces"},
	}
	for i, test := range tests {
		p, v := SplitHeaderOnSpace(test.in)
		assert.Check(t, cmp.Equal(test.p, p), "prefix mismatch on test index %d", i)
		assert.Check(t, cmp.Equal(test.v, v), "value mismatch on test index %d", i)
	}
}
