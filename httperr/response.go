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
)

// Response is an alias for http.Response that implements
// the error interface. Example:
//
//   resp, err := http.Get("http://www.example.com")
//   if err != nil {
//   	return err
//   }
//   if resp.StatusCode != http.StatusOK {
//   	return httperr.Response(*resp)
//   }
//   // ...
//
type Response http.Response

func (re Response) Error() string {
	msg := re.Header.Get("X-Error-Message")
	if msg != "" {
		msg = ": " + msg
	}
	return re.Status + msg
}
