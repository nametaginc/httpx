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

import "net/http"

// Code returns an http Error with the specified code, but without a wrapped error.
func Code(statusCode int) error {
	return Wrapf(statusCode, "%s", http.StatusText(statusCode))
}

var (
	// BadRequest is an error that returns a generic status 400 error
	BadRequest = Code(400)

	// Unauthorized is an error that returns a generic status 401 error
	Unauthorized = Code(401)

	// PaymentRequired is an error that returns a generic status 402 error
	PaymentRequired = Code(402)

	// Forbidden is an error that returns a generic status 403 error
	Forbidden = Code(403)

	// NotFound is an error that returns a generic status 404 error
	NotFound = Code(404)

	// MethodNotAllowed is an error that returns a generic status 405 error
	MethodNotAllowed = Code(405)

	// NotAcceptable is an error that returns a generic status 406 error
	NotAcceptable = Code(406)

	// ProxyAuthRequired is an error that returns a generic status 407 error
	ProxyAuthRequired = Code(407)

	// RequestTimeout is an error that returns a generic status 408 error
	RequestTimeout = Code(408)

	// Conflict is an error that returns a generic status 409 error
	Conflict = Code(409)

	// Gone is an error that returns a generic status 410 error
	Gone = Code(410)

	// LengthRequired is an error that returns a generic status 411 error
	LengthRequired = Code(411)

	// PreconditionFailed is an error that returns a generic status 412 error
	PreconditionFailed = Code(412)

	// RequestEntityTooLarge is an error that returns a generic status 413 error
	RequestEntityTooLarge = Code(413)

	// RequestURITooLong is an error that returns a generic status 414 error
	RequestURITooLong = Code(414)

	// UnsupportedMediaType is an error that returns a generic status 415 error
	UnsupportedMediaType = Code(415)

	// RequestedRangeNotSatisfiable is an error that returns a generic status 416 error
	RequestedRangeNotSatisfiable = Code(416)

	// ExpectationFailed is an error that returns a generic status 417 error
	ExpectationFailed = Code(417)

	// Teapot is an error that returns a generic status 418 error
	Teapot = Code(418)

	// InternalServerError is an error that returns a generic status 500 error
	InternalServerError = Code(500)

	// NotImplemented is an error that returns a generic status 501 error
	NotImplemented = Code(501)

	// BadGateway is an error that returns a generic status 502 error
	BadGateway = Code(502)

	// ServiceUnavailable is an error that returns a generic status 503 error
	ServiceUnavailable = Code(503)

	// GatewayTimeout is an error that returns a generic status 504 error
	GatewayTimeout = Code(504)

	// HTTPVersionNotSupported is an error that returns a generic status 505 error
	HTTPVersionNotSupported = Code(505)
)
