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
)

// BasicAuthTransport is an http.RoundTripper that adds a fixed basic auth
// to the request.
type BasicAuthTransport struct {
	Next     http.RoundTripper
	Username string
	Password string
}

// RoundTrip implements http.RoundTripper.
func (t BasicAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.Username, t.Password)
	return t.Next.RoundTrip(r)
}
