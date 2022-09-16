// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"net/http"

	"github.com/vitelabs/vite-portal/shared/pkg/crypto"
)

type jwtHandler struct {
	inner *crypto.JWTHandler
	next  http.Handler
}

// newJWTHandler creates a http.Handler with jwt authentication support.
func newJWTHandler(secret []byte, next http.Handler) http.Handler {
	return &jwtHandler{
		inner: crypto.NewDefaultJWTHandler(secret),
		next:  next,
	}
}

// ServeHTTP implements http.Handler
func (handler *jwtHandler) ServeHTTP(out http.ResponseWriter, r *http.Request) {
	token, err := handler.inner.Extract(r.Header)
	if err != nil {
		http.Error(out, err.Error(), http.StatusForbidden)
		return
	}
	_, err = handler.inner.Validate(token)
	if err != nil {
		http.Error(out, err.Error(), http.StatusForbidden)
		return
	}
	handler.next.ServeHTTP(out, r)
}
