// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A resolver that returns method values from a type switch keeps those
// methods and nothing else.

package main

import (
	"io"
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

type hello struct{}

func (hello) Resolve(receiver any, name string) any {
	if name != "Hello" {
		return nil
	}
	switch r := receiver.(type) {
	case A:
		return r.Hello
	case B:
		return r.Hello
	}
	return nil
}

func main() {
	t := template.Must(template.New("t").Parse("{{.}}"))
	t.ExecuteWithResolver(io.Discard, A(1), hello{})
	t.ExecuteWithResolver(io.Discard, B(1), hello{})
}
