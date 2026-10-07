// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build templatenomethods

package template

import (
	"strings"
	"testing"
)

// mapWithMethod has a method whose name is also one of its keys.
type mapWithMethod map[string]string

func (mapWithMethod) Method0() string { return "method" }

func TestExecuteNoMethods(t *testing.T) {
	const noMethod = "can't evaluate field Method0"
	tests := []struct {
		name, text string
		data       any
		want       string // the output, or the error text if wantErr is set
		wantErr    bool
	}{
		{"field", "{{.X}}", tVal, "x", false},
		{"method", "{{.Method0}}", tVal, noMethod, true},
		{"map key named like a method", "{{.Method0}}", mapWithMethod{"Method0": "key"}, "key", false},
		{"nested field", `{{define "inner"}}{{.X}}{{end}}{{template "inner" .}}`, tVal, "x", false},
		{"nested method", `{{define "inner"}}{{.Method0}}{{end}}{{template "inner" .}}`, tVal, noMethod, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpl := Must(New(test.name).Parse(test.text))
			var b strings.Builder
			err := tmpl.Execute(&b, test.data)
			if test.wantErr {
				if err == nil || !strings.Contains(err.Error(), test.want) {
					t.Fatalf("Execute = %q, %v; want error containing %q", b.String(), err, test.want)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if got := b.String(); got != test.want {
				t.Fatalf("Execute = %q; want %q", got, test.want)
			}
		})
	}
}
