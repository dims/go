// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The html/template form of tmplmethodmap_nil.

package main

import (
	"html/template"
	"io"
)

type T int

func (T) M() {}

func main() {
	template.Must(template.New("t").Parse("{{.}}")).ExecuteWithMethods(io.Discard, T(1), nil)
}
