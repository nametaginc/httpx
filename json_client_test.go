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
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"

	"github.com/nametaginc/httpx/httperr"
)

type APIError struct {
	Code    int
	Message string `json:"message"`
}

func (e *APIError) Error() string { return fmt.Sprintf("%s (%d)", e.Message, e.Code) }

type RequestBody struct {
	Foo string
}

type ResponseBody struct {
	Bar string
}

type FakeServer func(r *http.Request) (*http.Response, error)

func (f FakeServer) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type CannotMarshal struct{}

func (CannotMarshal) MarshalJSON() ([]byte, error) {
	return nil, errors.New("cannot frob the grob")
}

func TestJSONClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, r.Context(), ctx)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "https://api.example.com/foo", r.URL.String())
				assert.Equal(t, r.Header.Get("Content-type"), "application/json")
				assert.Equal(t, r.Header.Get("Accept"), "application/json")
				reqBody, err := ioutil.ReadAll(r.Body)
				assert.Check(t, err)
				assert.Equal(t, `{"Foo":"foo"}`, string(reqBody))

				resp := &http.Response{}
				resp.Body = ioutil.NopCloser(strings.NewReader(`{"Bar": "baz"}`))
				return resp, nil
			}),
		}}
		c.OnError = HandleJSONError(func() error { panic("not reached") })

		req := RequestBody{Foo: "foo"}
		resp := ResponseBody{}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", req, &resp)
		assert.Check(t, err)
		assert.DeepEqual(t, resp, ResponseBody{Bar: "baz"})
	})

	t.Run("parse-error", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, r.Context(), ctx)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "https://api.example.com/foo", r.URL.String())
				reqBody, err := ioutil.ReadAll(r.Body)
				assert.Check(t, err)
				assert.Equal(t, `{"Foo":"foo"}`, string(reqBody))

				resp := &http.Response{}
				resp.StatusCode = http.StatusTeapot
				resp.Body = ioutil.NopCloser(strings.NewReader(`{"Code": 1337, "Message": "good hacker. have cookie"}`))
				return resp, nil
			}),
		}}
		c.OnError = HandleJSONError(func() error { return &APIError{} })

		req := RequestBody{Foo: "foo"}
		resp := ResponseBody{}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", req, &resp)
		assert.DeepEqual(t, err.(*APIError), &APIError{
			Code:    1337,
			Message: "good hacker. have cookie",
		})
		assert.DeepEqual(t, resp, ResponseBody{})
	})

	t.Run("no request body", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, r.Context(), ctx)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, r.Header.Get("Content-type"), "")
				assert.Equal(t, r.Header.Get("Accept"), "application/json")
				assert.Equal(t, "https://api.example.com/foo", r.URL.String())
				assert.Check(t, is.Nil(r.Body))

				resp := &http.Response{}
				resp.Body = ioutil.NopCloser(strings.NewReader(`{"Bar": "baz"}`))
				return resp, nil
			}),
		}}
		c.OnError = HandleJSONError(func() error { panic("not reached") })

		resp := ResponseBody{}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", nil, &resp)
		assert.Check(t, err)
		assert.DeepEqual(t, resp, ResponseBody{Bar: "baz"})
	})

	t.Run("no response body", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, r.Context(), ctx)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "https://api.example.com/foo", r.URL.String())
				assert.Equal(t, r.Header.Get("Content-type"), "application/json")
				assert.Equal(t, r.Header.Get("Accept"), "")
				reqBody, err := ioutil.ReadAll(r.Body)
				assert.Check(t, err)
				assert.Equal(t, `{"Foo":"foo"}`, string(reqBody))

				// the client isn't requesting a response, but we can still send one anyway
				resp := &http.Response{}
				resp.Body = ioutil.NopCloser(strings.NewReader(`{"Bar": "baz"}`))
				return resp, nil
			}),
		}}
		c.OnError = HandleJSONError(func() error { panic("not reached") })

		req := RequestBody{Foo: "foo"}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", req, nil)
		assert.Check(t, err)
	})

	t.Run("error cannot be parsed", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, r.Context(), ctx)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "https://api.example.com/foo", r.URL.String())
				reqBody, err := ioutil.ReadAll(r.Body)
				assert.Check(t, err)
				assert.Equal(t, `{"Foo":"foo"}`, string(reqBody))

				resp := &http.Response{}
				resp.StatusCode = http.StatusTeapot
				resp.Body = ioutil.NopCloser(strings.NewReader(`{invalid json`))
				return resp, nil
			}),
		}}
		c.OnError = HandleJSONError(func() error { return &APIError{} })

		req := RequestBody{Foo: "foo"}
		resp := ResponseBody{}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", req, &resp)
		assert.Equal(t, err.(httperr.Response).StatusCode, http.StatusTeapot)
		body, err := ioutil.ReadAll(err.(httperr.Response).Body)
		assert.Check(t, err)
		assert.Equal(t, string(body), "{invalid json")
	})

	t.Run("cannot marshal request", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				panic("not reached")
			}),
		}}
		c.OnError = HandleJSONError(func() error { panic("not reached") })

		req := CannotMarshal{}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", req, nil)
		assert.Error(t, err, "json: error calling MarshalJSON for type httpx.CannotMarshal: cannot frob the grob")
	})

	t.Run("cannot unmarshal response", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				resp := &http.Response{}
				resp.Body = ioutil.NopCloser(strings.NewReader(`{invalid json`))
				return resp, nil
			}),
		}}
		c.OnError = HandleJSONError(func() error { panic("not reached") })

		resp := ResponseBody{}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", nil, &resp)
		assert.Error(t, err, "invalid character 'i' looking for beginning of object key string")
	})

	t.Run("default error handler", func(t *testing.T) {
		ctx := context.Background()
		c := JSONClient{Client: &http.Client{
			Transport: FakeServer(func(r *http.Request) (*http.Response, error) {
				resp := &http.Response{}
				resp.StatusCode = http.StatusTeapot
				resp.Body = ioutil.NopCloser(strings.NewReader(`{"Code": 1337, "Message": "good hacker. have cookie"}`))
				return resp, nil
			}),
		}}

		resp := ResponseBody{}
		err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", nil, &resp)
		assert.Equal(t, err.(httperr.Response).StatusCode, http.StatusTeapot)
		body, err := ioutil.ReadAll(err.(httperr.Response).Body)
		assert.Check(t, err)
		assert.Equal(t, string(body), `{"Code": 1337, "Message": "good hacker. have cookie"}`)
	})
}
