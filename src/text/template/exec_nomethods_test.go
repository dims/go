// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"strings"
	"testing"
)

type noMethodsData struct{ Field string }

func (noMethodsData) Method() string { return "method" }

func TestExecuteWithoutMethods(t *testing.T) {
	tests := []struct {
		name, text string
		data       any
		want, err  string
	}{
		{"field", "{{.Field}}", noMethodsData{Field: "field"}, "field", ""},
		{"method", "{{.Method}}", noMethodsData{}, "", "can't evaluate field Method"},
		{"map key named like a method", "{{.Method}}", map[string]string{"Method": "key"}, "key", ""},
		{"nested template", `{{define "x"}}{{.Method}}{{end}}{{template "x" .}}`, noMethodsData{}, "", "can't evaluate field Method"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			err := Must(New(tc.name).Parse(tc.text)).ExecuteWithoutMethods(&b, tc.data)
			if tc.err == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.err != "" && (err == nil || !strings.Contains(err.Error(), tc.err)) {
				t.Fatalf("error %v, want one containing %q", err, tc.err)
			}
			if got := b.String(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}

	// Lookup followed by ExecuteWithoutMethods is the named-template form.
	root := Must(New("root").Parse(`{{define "x"}}{{.Method}}{{end}}`))
	var b strings.Builder
	if err := root.Lookup("x").ExecuteWithoutMethods(&b, noMethodsData{}); err == nil || !strings.Contains(err.Error(), "can't evaluate field Method") {
		t.Fatalf("Lookup: %v, want an error", err)
	}
	// Execute on the same template and data still calls the method.
	if err := root.Lookup("x").Execute(&b, noMethodsData{}); err != nil || b.String() != "method" {
		t.Fatalf("Execute: got %q, %v; want the method result", b.String(), err)
	}
}
