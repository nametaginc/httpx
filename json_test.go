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
	"net/http"
	"testing"

	"gotest.tools/assert"
)

type TestInputType struct{}
type TestOutputType struct{}

func TestHttpHandlerTypeSignaturesArePanicFree(t *testing.T) {
	requestArgReturnsOutputPtrError := func(r *http.Request, in TestInputType) (*TestOutputType, error) {
		return nil, nil
	}
	assert.Check(t, JSONHandler(requestArgReturnsOutputPtrError) != nil)
	requestPtrArgReturnsOutputPtrError := func(r *http.Request, in *TestInputType) (*TestOutputType, error) {
		return nil, nil
	}
	assert.Check(t, JSONHandler(requestPtrArgReturnsOutputPtrError) != nil)
	requestArgReturnsError := func(r *http.Request, in TestInputType) error {
		return nil
	}
	assert.Check(t, JSONHandler(requestArgReturnsError) != nil)
	requestPtrArgReturnsError := func(r *http.Request, in *TestInputType) error {
		return nil
	}
	assert.Check(t, JSONHandler(requestPtrArgReturnsError) != nil)
	requestReturnsError := func(r *http.Request) error {
		return nil
	}
	assert.Check(t, JSONHandler(requestReturnsError) != nil)
	requestOnly := func(r *http.Request) {
		return
	}
	assert.Check(t, JSONHandler(requestOnly) != nil)
}
