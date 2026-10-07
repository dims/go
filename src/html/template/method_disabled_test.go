// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build templatenomethods

package template

// The templatenomethods build tag disables method calls in text/template.
const methodsEnabled = false
