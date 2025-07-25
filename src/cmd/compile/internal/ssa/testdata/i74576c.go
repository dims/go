// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"runtime"
)

func main() {
	s := S{1, 1}
	runtime.Breakpoint()
	sink = s
}

type S struct{ a, b uint64 }

var sink any
