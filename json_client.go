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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nametaginc/httpx/httperr"
)

// JSONClient extends http.Client which adds a method that makes handling
// JSON requests / responses easier.
//
// Example:
//
//    type APIError struct {
//      Code int `json:"code"`
//      Message string `json:"message"`
//    }
//    func (e APIError) Error() string { return fmt.Sprintf("%s (%d)", e.Message, e.String)
//
//    c := JSONClient{Client: http.DefaultClient}
//    c.OnError = HandleJSONError(func() error { return &APIError{} })
//
//    req := struct { Foo string }{}
//    resp := string { Bar string }{}
//    if err := c.DoJSON(ctx, "POST", "https://api.example.com/foo", req, &resp); err != nil {
//      return err
//    }
//
type JSONClient struct {
	*http.Client
	OnError func(r *http.Response) error
}

// HandleJSONError is a function that can be assigned to JSONClient.OnError to handle
// error bodies that are returned by an API
func HandleJSONError(respBodyFactory func() error) func(r *http.Response) error {
	return func(r *http.Response) error {
		respBodyBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		respBody := respBodyFactory()
		if err := json.Unmarshal(respBodyBuf, respBody); err != nil {
			r.Body = ioutil.NopCloser(bytes.NewReader(respBodyBuf))
			return httperr.Response(*r)
		}
		return respBody
	}
}

// NewRequest returns a new request having the requestBody as the HTTP request body serialized in JSON format.
func (c JSONClient) NewRequest(ctx context.Context, method string, uri string, requestBody interface{}) (*http.Request, error) {
	var body io.Reader
	if requestBody != nil {
		bodyBuf, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bodyBuf)
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	if req.Body != nil {
		req.Header.Add("Content-type", "application/json")
	}
	return req, nil
}

// HandleResponse handles an HTTP response containing a JSON object. If responseBody is provided, then
// the response is unmarshalled into it. If the HTTP status code is >= 400, then OnError is invoked if
// provided, otherwise an httperr.Response error is returned.
func (c JSONClient) HandleResponse(resp *http.Response, responseBody interface{}) error {
	if resp.StatusCode >= 400 {
		if c.OnError != nil {
			return c.OnError(resp)
		}
		return httperr.Response(*resp)
	}

	if responseBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
			return err
		}
	}

	return nil
}

// DoJSON performs an HTTP request. If request is provided, it marshalled into the request body. If response is
// provided, the response body is marshalled into it. If the server returns an error, then the response is passed
// through OnError, or if OnError is not provided, an httperr.Response error is returned.
func (c JSONClient) DoJSON(ctx context.Context, method string, uri string, request interface{}, response interface{}) error {
	httpReq, err := c.NewRequest(ctx, method, uri, request)
	if err != nil {
		return err
	}
	if response != nil {
		httpReq.Header.Set("Accept", "application/json")
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return err
	}
	return c.HandleResponse(httpResp, response)
}
