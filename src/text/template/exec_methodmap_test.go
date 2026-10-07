// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"errors"
	"strings"
	"testing"
)

type methodMapData struct{ Field string }

func (methodMapData) Method() string { return "method" }

func TestExecuteWithMethods(t *testing.T) {
	methods := FuncMap{
		"Value":   func(d methodMapData) string { return "value " + d.Field },
		"Pointer": func(d *methodMapData) string { return "pointer " + d.Field },
		"Add":     func(d methodMapData, n int) int { return len(d.Field) + n },
		"Fail":    func(methodMapData) (string, error) { return "", errors.New("failed") },
		"Method":  func(methodMapData) string { return "mapped" },
	}
	tests := []struct {
		name, text string
		data       any
		methods    FuncMap
		want, err  string
	}{
		{"value receiver", "{{.Value}}", methodMapData{"x"}, methods, "value x", ""},
		{"value receiver from pointer", "{{.Value}}", &methodMapData{"x"}, methods, "value x", ""},
		{"pointer receiver", "{{.Pointer}}", &methodMapData{"x"}, methods, "pointer x", ""},
		{"pointer receiver, unaddressable data", "{{.Pointer}}", methodMapData{"x"}, methods, "", "wrong receiver type for method Pointer; expected *template.methodMapData; got template.methodMapData"},
		{"arguments", "{{.Add 2}}", methodMapData{"x"}, methods, "3", ""},
		{"pipeline argument", "{{2 | .Add}}", methodMapData{"x"}, methods, "3", ""},
		{"too few arguments", "{{.Add}}", methodMapData{"x"}, methods, "", "wrong number of args for Add: want 1 got 0"},
		{"error result", "{{.Fail}}", methodMapData{}, methods, "", "error calling Fail: failed"},
		{"map wins over a real method", "{{.Method}}", methodMapData{}, methods, "mapped", ""},
		{"field", "{{.Field}}", methodMapData{"x"}, methods, "x", ""},
		{"map key", "{{.Field}}", map[string]string{"Field": "key"}, methods, "key", ""},
		{"name not in the map", "{{.Other}}", methodMapData{}, methods, "", "can't evaluate field Other"},
		{"nil map, field", "{{.Field}}", methodMapData{"x"}, nil, "x", ""},
		{"nil map, method", "{{.Method}}", methodMapData{}, nil, "", "can't evaluate field Method"},
		{"variable", "{{$x := .}}{{$x.Value}}", methodMapData{"x"}, methods, "value x", ""},
		{"nested template", `{{define "x"}}{{.Value}}{{end}}{{template "x" .}}`, methodMapData{"x"}, methods, "value x", ""},
		{"nested template, nil map", `{{define "x"}}{{.Method}}{{end}}{{template "x" .}}`, methodMapData{}, nil, "", "can't evaluate field Method"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			err := Must(New(tc.name).Parse(tc.text)).ExecuteWithMethods(&b, tc.data, tc.methods)
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
}

// A bad map is rejected before execution writes anything.
func TestExecuteWithMethodsBadMap(t *testing.T) {
	tests := []struct {
		name string
		fn   any
		err  string
	}{
		{"not a function", 42, `method "M": value is not a function`},
		{"nil", nil, `method "M": value is not a function`},
		{"no receiver", func() string { return "" }, `method "M": function has no receiver parameter`},
		{"no result", func(methodMapData) {}, `method "M": function M has 0 return values`},
		{"bad second result", func(methodMapData) (int, int) { return 0, 0 }, `method "M": invalid function signature for M: second return value should be error`},
	}
	tmpl := Must(New("t").Parse("text{{.M}}"))
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			err := tmpl.ExecuteWithMethods(&b, methodMapData{}, FuncMap{"M": tc.fn})
			if err == nil || !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("error %v, want one containing %q", err, tc.err)
			}
			if b.Len() != 0 {
				t.Fatalf("wrote %q before rejecting the map", b.String())
			}
		})
	}
}

// Lookup followed by ExecuteWithMethods is the named-template form, and
// Execute still resolves real methods.
func TestExecuteWithMethodsLookup(t *testing.T) {
	root := Must(New("root").Parse(`{{define "x"}}{{.Method}}{{end}}`))
	var b strings.Builder
	err := root.Lookup("x").ExecuteWithMethods(&b, methodMapData{}, FuncMap{"Method": func(methodMapData) string { return "mapped" }})
	if err != nil || b.String() != "mapped" {
		t.Fatalf("ExecuteWithMethods: got %q, %v; want \"mapped\"", b.String(), err)
	}
	b.Reset()
	if err := root.Lookup("x").Execute(&b, methodMapData{}); err != nil || b.String() != "method" {
		t.Fatalf("Execute: got %q, %v; want \"method\"", b.String(), err)
	}
}
