// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"strings"
	"testing"
)

type resolverData struct{ Field string }

func (resolverData) Method() string { return "<method>" }

type methodOnly struct{}

func (methodOnly) Resolve(receiver any, name string) any {
	if r, ok := receiver.(resolverData); ok && name == "Method" {
		return r.Method
	}
	return nil
}

func TestExecuteWithResolver(t *testing.T) {
	var b strings.Builder
	err := Must(New("t").Parse("{{.Field}}")).ExecuteWithResolver(&b, resolverData{Field: "<field>"}, nil)
	if err != nil || b.String() != "&lt;field&gt;" {
		t.Fatalf("field with nil resolver: got %q, %v; want the escaped field", b.String(), err)
	}
	err = Must(New("t").Parse("{{.Method}}")).ExecuteWithResolver(&b, resolverData{}, nil)
	if err == nil || !strings.Contains(err.Error(), "can't evaluate field Method") {
		t.Fatalf("method with nil resolver: %v, want an error", err)
	}
	b.Reset()
	err = Must(New("t").Parse("{{.Method}}")).ExecuteWithResolver(&b, resolverData{}, methodOnly{})
	if err != nil || b.String() != "&lt;method&gt;" {
		t.Fatalf("method with a resolver: got %q, %v; want the escaped method result", b.String(), err)
	}
	b.Reset()
	root := Must(New("root").Parse(`{{define "x"}}{{.Field}}{{end}}`))
	if err := root.Lookup("x").ExecuteWithResolver(&b, resolverData{Field: "<f>"}, nil); err != nil || b.String() != "&lt;f&gt;" {
		t.Fatalf("Lookup: got %q, %v; want the escaped field", b.String(), err)
	}
}
