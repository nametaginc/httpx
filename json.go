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
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/nametaginc/httpx/httperr"
)

// JSONHandler returns an http handler that accepts JSON as input and emits JSON as output. The input and output
// are serialized.
//
// Example usage:
//
//   mux.Post("/someendpoint", JSONHandler(func(r *http.Request, in InputType) (*OutputType, error) {
//      /* implementation */
//   }))
//
// InputType and OutputType must be structs.
//
// The function must have one of the following signatures:
//
//    func(r *http.Request, in InputType) (*OutputType, error)
//    func(r *http.Request, in *InputType) (*OutputType, error)
//    func(r *http.Request) (*OutputType, error)
//    func(r *http.Request, in InputType) (error)
//    func(r *http.Request, in *InputType) (error)
//    func(r *http.Request) (error)
//    func(r *http.Request)
//
func JSONHandler(f interface{}) http.Handler {
	fval := reflect.ValueOf(f)
	ftyp := fval.Type()

	if ftyp.NumIn() != 1 && ftyp.NumIn() != 2 {
		panic("function must take 1 or 2 two arguments")
	}
	if ftyp.In(0).Kind() != reflect.Ptr || ftyp.In(0).Elem() != reflect.TypeOf(http.Request{}) {
		panic("first argument must be *http.Request")
	}
	if ftyp.NumIn() == 2 {
		arg := ftyp.In(1)
		kind := arg.Kind()
		if kind == reflect.Ptr && arg.Elem().Kind() != reflect.Struct {
			panic("second argument must be a pointer to a struct type, got pointer to non-struct")
		}
		if kind != reflect.Ptr && kind != reflect.Struct {
			panic("second argument must be a struct or pointer to stuct")
		}
	}

	if ftyp.NumOut() == 1 {
		errTyp := ftyp.Out(0)
		if errTyp.String() != "error" { // TODO(ross): do a better type comparison
			panic("function must return (error)")
		}
	} else if ftyp.NumOut() == 2 {
		outTyp := ftyp.Out(0)
		if outTyp.Kind() != reflect.Ptr || outTyp.Elem().Kind() != reflect.Struct {
			panic("function must return (*SomeStruct, error)")
		}

		errTyp := ftyp.Out(1)
		if errTyp.String() != "error" { // TODO(ross): do a better type comparison
			panic("function must return (*SomeStruct, error)")
		}
	}

	return httperr.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var out []reflect.Value
		if ftyp.NumIn() == 2 {
			reqBody := reflect.New(ftyp.In(1))
			if err := json.NewDecoder(r.Body).Decode(reqBody.Interface()); err != nil {
				return httperr.Public(http.StatusBadRequest, err)
			}

			out = fval.Call([]reflect.Value{reflect.ValueOf(r), reqBody.Elem()})
		} else if ftyp.NumIn() == 1 {
			out = fval.Call([]reflect.Value{reflect.ValueOf(r)})
		}

		var err error
		if ftyp.NumOut() == 1 {
			if e := out[0].Interface(); e != nil {
				err = e.(error)
			}
		} else if ftyp.NumOut() == 2 {
			if e := out[1].Interface(); e != nil {
				err = e.(error)
			}
		}
		if err != nil {
			return err
		}

		if ftyp.NumOut() == 1 {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}

		respBody := out[0].Interface()
		w.Header().Add("Content-type", "application/json")
		return json.NewEncoder(w).Encode(respBody)
	})
}
