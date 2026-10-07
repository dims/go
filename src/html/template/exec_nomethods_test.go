// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"strings"
	"testing"
)

type noMethodsData struct{ Field string }

func (noMethodsData) Method() string { return "<method>" }

func TestExecuteWithoutMethods(t *testing.T) {
	var b strings.Builder
	err := Must(New("t").Parse("{{.Field}}")).ExecuteWithoutMethods(&b, noMethodsData{Field: "<field>"})
	if err != nil || b.String() != "&lt;field&gt;" {
		t.Fatalf("field: got %q, %v; want the escaped field", b.String(), err)
	}
	err = Must(New("t").Parse("{{.Method}}")).ExecuteWithoutMethods(&b, noMethodsData{})
	if err == nil || !strings.Contains(err.Error(), "can't evaluate field Method") {
		t.Fatalf("method: %v, want an error", err)
	}
	b.Reset()
	root := Must(New("root").Parse(`{{define "x"}}{{.Field}}{{end}}`))
	if err := root.Lookup("x").ExecuteWithoutMethods(&b, noMethodsData{Field: "<f>"}); err != nil || b.String() != "&lt;f&gt;" {
		t.Fatalf("Lookup: got %q, %v; want the escaped field", b.String(), err)
	}
}
