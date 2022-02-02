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

// ErrorTranslatorFunc is a function that translates errors
// before they are returned by the http client.
type ErrorTranslatorFunc func(err error) error

type errorTranslatorFuncHolder struct {
	F ErrorTranslatorFunc
}

// RegisterErrorTranslator adds a new error translation function to the
// error translation stack. Your function will be called with each
// error before writing the response.
//
// Note: this function is not go-routine safe. You should probably call
// it from an init() function or similar.
//
// Example:
//
//   init() {
//     RegisterErrorTranslator(func(err error) error {
//       if (err == context.DeadlineExceeded) {
//         return Error{
//           PrivateError: err,
//           Code: http.StatusRequestTimeout,
//         }
//       }
//     })
//   }
//
func RegisterErrorTranslator(f ErrorTranslatorFunc) {
	holder := &errorTranslatorFuncHolder{F: f}
	errorTranslators = append([]*errorTranslatorFuncHolder{holder}, errorTranslators...)
}

var errorTranslators []*errorTranslatorFuncHolder

// TranslateError invokes all the registered ErrorTranslatorFunc functions on
// err and returns the result.
func TranslateError(err error) error {
	for _, et := range errorTranslators {
		err = et.F(err)
	}
	return err
}
