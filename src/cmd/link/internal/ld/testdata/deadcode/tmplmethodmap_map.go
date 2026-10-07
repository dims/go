// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A method FuncMap keeps the functions it references and nothing else.

package main

import (
	"io"
	"text/template"
)

type T int

func (T) M() {}
func (T) N() {}

func main() {
	t := template.Must(template.New("t").Parse("{{.}}"))
	t.ExecuteWithMethods(io.Discard, T(1), template.FuncMap{"M": T.M})
}
