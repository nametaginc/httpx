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
	"net/http"
	"strconv"
	"time"
)

var Now = time.Now

// TooManyRequests is an error that returns an HTTP status 429 response when a client
// has violated a rate limit restriction.
type TooManyRequests struct {
	RetryAfter *time.Time
}

var _ ResponseWriter = TooManyRequests{}
var _ StatusCoder = TooManyRequests{}

func (e TooManyRequests) Error() string {
	return http.StatusText(e.StatusCode())
}

// WriteResponse implements ResponseWriter
func (e TooManyRequests) WriteResponse(w http.ResponseWriter, r *http.Request) {
	if e.RetryAfter != nil {
		// Retry-After can be either an HTTP date or a number of seconds. Because clock-skew
		// is still a thing in 2021, we use the `delay-seconds` flavor.
		//
		// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Retry-After
		retryAfter := e.RetryAfter.Sub(Now()).Seconds()
		if retryAfter < 0 {
			// delay-seconds is a "non-negative decimal integer, representing time in seconds"
			// ref: https://httpwg.org/specs/rfc7231.html#header.retry-after
			retryAfter = 0
		}

		w.Header().Add("Retry-After", strconv.Itoa(int(retryAfter)))
	}
	http.Error(w, http.StatusText(e.StatusCode()), e.StatusCode())
}

// StatusCode implements StatusCoder
func (e TooManyRequests) StatusCode() int {
	return http.StatusTooManyRequests
}
