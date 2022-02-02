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
	"errors"
	"net/http"

	pkgerrors "github.com/pkg/errors"
)

// ResponseWriter is an interface implemented by errors that format themselves as an HTTP response.
type ResponseWriter interface {
	error
	WriteResponse(w http.ResponseWriter, r *http.Request)
}

// Write emits an error message to w. If err implements ResponseWriter, then it is used to render the response. If err
// implements StatusCoder, then it is used to set the status code. If err implements IsPublicer, and IsPublic() returns
// true, then the text of the error is included in the response, otherwise a generic message, e.g. "Bad Request" is
// included in the response.
func Write(w http.ResponseWriter, r *http.Request, err error) {
	var rw ResponseWriter
	if errors.As(err, &rw) {
		rw.WriteResponse(w, r)
		return
	}

	statusCode := StatusCode(err)
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	if IsPublic(err) {
		errText := err.Error()

		if rootErr := pkgerrors.Cause(err); rootErr != nil {
			errText = rootErr.Error()
		}

		w.Header().Add("X-Error-Message", errText)
	}
	http.Error(w, http.StatusText(statusCode), statusCode)
}
