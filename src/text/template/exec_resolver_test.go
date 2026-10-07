// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"reflect"
	"strings"
	"testing"
)

type resolverData struct{ Field string }

func (resolverData) Method() string { return "method" }
func (resolverData) Other() string  { return "other" }

// onlyMethod resolves Method and nothing else, by type switch.
type onlyMethod struct{}

func (onlyMethod) Resolve(receiver any, name string) any {
	if name != "Method" {
		return nil
	}
	switch r := receiver.(type) {
	case resolverData:
		return r.Method
	case *resolverData:
		return r.Method
	}
	return nil
}

// byReflection resolves every method, as Execute does.
type byReflection struct{}

func (byReflection) Resolve(receiver any, name string) any {
	m := reflect.ValueOf(receiver).MethodByName(name)
	if !m.IsValid() {
		return nil
	}
	return m.Interface()
}

// panicky panics, as a careless resolver might.
type panicky struct{}

func (panicky) Resolve(any, string) any { panic("boom") }

// constant returns the same value for every name.
type constant struct{ v any }

func (c constant) Resolve(any, string) any { return c.v }

// shim maps Method on a resolverData to a function that is not a method
// of the type and takes an argument from the template.
type shim struct{}

func (shim) Resolve(receiver any, name string) any {
	if r, ok := receiver.(resolverData); ok && name == "Method" {
		return func(suffix string) string { return r.Field + suffix }
	}
	return nil
}

func TestExecuteWithResolver(t *testing.T) {
	tests := []struct {
		name, text string
		data       any
		resolver   MethodResolver
		want, err  string
	}{
		{"nil field", "{{.Field}}", resolverData{Field: "field"}, nil, "field", ""},
		{"nil method", "{{.Method}}", resolverData{}, nil, "", "can't evaluate field Method"},
		{"nil map key named like a method", "{{.Method}}", map[string]string{"Method": "key"}, nil, "key", ""},
		{"nil missing map key", "{{.Missing}}", map[string]string{}, nil, "<no value>", ""},
		{"nil in nested template", `{{define "x"}}{{.Method}}{{end}}{{template "x" .}}`, resolverData{}, nil, "", "can't evaluate field Method"},
		{"resolver allows one name", "{{.Method}}", resolverData{}, onlyMethod{}, "method", ""},
		{"resolver sees a pointer receiver", "{{.Method}}", &resolverData{}, onlyMethod{}, "method", ""},
		{"resolver denies another", "{{.Other}}", resolverData{}, onlyMethod{}, "", "can't evaluate field Other"},
		{"resolver in nested template", `{{define "x"}}{{.Method}}{{end}}{{template "x" .}}`, resolverData{}, onlyMethod{}, "method", ""},
		{"resolver by reflection", "{{.Method}} {{.Other}}", resolverData{}, byReflection{}, "method other", ""},
		{"resolver returning nil for everything", "{{.Method}}", resolverData{}, constant{nil}, "", "can't evaluate field Method"},
		{"resolver panics", "{{.Method}}", resolverData{}, panicky{}, "", "error calling method resolver for Method: boom"},
		{"resolver binds a method on a nil pointer", "{{.Method}}", (*resolverData)(nil), onlyMethod{}, "", "error calling method resolver for Method"},
		{"resolver shims a method with arguments", `{{.Method "!"}}`, resolverData{Field: "field"}, shim{}, "field!", ""},
		{"resolver returns a non-function", "{{.Method}}", resolverData{}, constant{42}, "", "returned int for Method, not a function"},
		{"resolver returns a bad signature", "{{.Method}}", resolverData{}, constant{func() (int, int) { return 1, 2 }}, "", "invalid function signature for Method"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			err := Must(New(tc.name).Parse(tc.text)).ExecuteWithResolver(&b, tc.data, tc.resolver)
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

// The missingkey option applies as it does to Execute.
func TestExecuteWithResolverMissingKey(t *testing.T) {
	var b strings.Builder
	err := Must(New("m").Option("missingkey=error").Parse("{{.Missing}}")).ExecuteWithResolver(&b, map[string]string{}, nil)
	if err == nil || !strings.Contains(err.Error(), `map has no entry for key "Missing"`) {
		t.Fatalf("missingkey=error: %v", err)
	}
}

// Lookup followed by ExecuteWithResolver is the named-template form.
func TestExecuteWithResolverLookup(t *testing.T) {
	root := Must(New("root").Parse(`{{define "x"}}{{.Method}}{{end}}`))
	var b strings.Builder
	if err := root.Lookup("x").ExecuteWithResolver(&b, resolverData{}, nil); err == nil || !strings.Contains(err.Error(), "can't evaluate field Method") {
		t.Fatalf("nil resolver: %v, want an error", err)
	}
	if err := root.Lookup("x").ExecuteWithResolver(&b, resolverData{}, onlyMethod{}); err != nil || b.String() != "method" {
		t.Fatalf("resolver: got %q, %v", b.String(), err)
	}
}
