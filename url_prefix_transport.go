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
	"net/url"
)

// URLPrefixTransport is an http.RoundTripper that prepends Server to each URL.
//
// e.g.
//
//   transport := URLPrefixTransport{Next: http.DefaultTransport, Server: "https://example.com/api/v1"}
//
type URLPrefixTransport struct {
	Next   http.RoundTripper
	Server string
}

// RoundTrip implements http.RoundTripper.
func (t URLPrefixTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var err error
	r.URL, err = url.Parse(t.Server + r.URL.String())
	if err != nil {
		return nil, err
	}
	return t.Next.RoundTrip(r)
}
