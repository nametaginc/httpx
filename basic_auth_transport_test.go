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

func TestBasicAuthTransport(t *testing.T) {
	fakeTransport := FakeServer(func(r *http.Request) (*http.Response, error) {
		username, password, ok := r.BasicAuth()
		assert.Check(t, ok)
		assert.Equal(t, "alice", username)
		assert.Equal(t, "hunter2", password)

		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(strings.NewReader(`Hello, World!`))
		return resp, nil
	})

	wrapped := BasicAuthTransport{
		Next:     fakeTransport,
		Username: "alice",
		Password: "hunter2",
	}

	client := http.Client{Transport: wrapped}
	resp, err := client.Get("dontcare")
	assert.Check(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.Check(t, err)
	assert.Equal(t, "Hello, World!", string(respBody))
}
