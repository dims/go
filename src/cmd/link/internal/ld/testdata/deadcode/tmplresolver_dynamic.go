// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A resolver that looks methods up by reflection keeps all exported methods,
// as Execute does.

package main

import (
	"io"
	"reflect"
	"text/template"
)

type A int
type B int

//go:noinline
func (A) Hello() {}

//go:noinline
func (A) Goodbye() {}

//go:noinline
func (B) Hello() {}

//go:noinline
func (B) Goodbye() {}

type reflective struct{}

func (reflective) Resolve(receiver any, name string) any {
	m := reflect.ValueOf(receiver).MethodByName(name)
	if !m.IsValid() {
		return nil
	}
	return m.Interface()
}

func main() {
	t := template.Must(template.New("t").Parse("{{.}}"))
	t.ExecuteWithResolver(io.Discard, A(1), reflective{})
	t.ExecuteWithResolver(io.Discard, B(1), reflective{})
}
