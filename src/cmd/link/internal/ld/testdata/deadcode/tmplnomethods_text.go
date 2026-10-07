// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A template executed with ExecuteWithoutMethods does not keep methods alive.

package main

import (
	"io"
	"text/template"
)

type T int

func (T) M() {}

func main() {
	template.Must(template.New("t").Parse("{{.}}")).ExecuteWithoutMethods(io.Discard, T(1))
}
