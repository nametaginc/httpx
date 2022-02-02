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
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestUrlPrefixTransport(t *testing.T) {
	fakeTransport := FakeServer(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, r.URL.String(), "https://api.example.com/v1/foo")
		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(strings.NewReader(`Hello, World!`))
		return resp, nil
	})

	wrapped := URLPrefixTransport{
		Next:   fakeTransport,
		Server: "https://api.example.com/v1",
	}

	client := http.Client{Transport: wrapped}
	resp, err := client.Get("/foo")
	assert.Check(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.Check(t, err)
	assert.Equal(t, "Hello, World!", string(respBody))
}
