// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !templatenomethods

package template

import "reflect"

const methodsEnabled = true

// methodByName returns the method of v named name, or the zero Value.
func methodByName(v reflect.Value, name string) reflect.Value { return v.MethodByName(name) }
