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
	"net/http"

	"goji.io/pattern"
)

// TestRequest returns a request suitable for use with an httpx.JSONHandler, optionally injecting
// goji variables
func TestRequest(ctx context.Context, vars ...string) *http.Request {
	r, err := http.NewRequest("DONTCARE", "dont://care", nil)
	if err != nil {
		panic(err)
	}
	if len(vars)%2 != 0 {
		panic("usage: TestRequest(ctx, [var, value]...")
	}
	for i := 0; i < len(vars); i += 2 {
		ctx = context.WithValue(ctx, pattern.Variable(vars[i]), vars[i+1])
	}
	return r.WithContext(ctx)
}
