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

// UserAgentTransport is an http.RoundTripper that adds a fixed user agent to the
// request if it is not already set.
type UserAgentTransport struct {
	Next      http.RoundTripper
	UserAgent string
}

// RoundTrip implements http.RoundTripper.
func (t UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Header.Get("User-Agent") == "" {
		r.Header.Set("User-Agent", t.UserAgent)
	}
	return t.Next.RoundTrip(r)
}
