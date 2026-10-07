// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build templatenomethods

package template

import "reflect"

const methodsEnabled = false

// methodByName never finds a method, so the linker can drop methods that
// no code calls directly.
func methodByName(reflect.Value, string) reflect.Value { return reflect.Value{} }
