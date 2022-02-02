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
	"strings"
)

// SplitHeaderOnSpace returns the prefix and remaining value of a header at the
// first set of consecutive whitespace.
func SplitHeaderOnSpace(h string) (prefix, value string) {
	fields := strings.Fields(h)
	if len(fields) == 2 {
		return fields[0], fields[1]
	}
	if len(fields) > 2 {
		prefixEnd := len(fields[0])
		valueOffset := strings.Index(h[prefixEnd:], fields[1])
		return fields[0], h[prefixEnd+valueOffset:]
	}
	return "", h
}

// SplitAuthorizationHeader splits the parts of an Authorization header of the form "Method Token OtherStuff"
func SplitAuthorizationHeader(r *http.Request) (kind string, value string) {
	kind, value = SplitHeaderOnSpace(r.Header.Get("Authorization"))
	return strings.ToLower(kind), value
}
