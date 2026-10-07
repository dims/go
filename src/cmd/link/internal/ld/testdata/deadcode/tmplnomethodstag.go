// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Executing a text/template reaches reflect.Value.MethodByName, which
// keeps every exported method live, unless the program is built with
// the templatenomethods tag.

package main

import (
	"io"
	"text/template"
)

type T int

func (T) M() {}

func main() {
	t := template.Must(template.New("t").Parse("{{.}}"))
	t.Execute(io.Discard, T(1))
}
