// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"strings"
	"testing"
)

type methodMapData struct{ Field string }

func (methodMapData) Method() string { return "<method>" }

func TestExecuteWithMethods(t *testing.T) {
	methods := FuncMap{"Mapped": func(d methodMapData) string { return "<" + d.Field + ">" }}
	tests := []struct {
		name, text string
		methods    FuncMap
		want, err  string
	}{
		{"mapped result is escaped", "{{.Mapped}}", methods, "&lt;x&gt;", ""},
		{"field is escaped", "{{.Field}}", nil, "x", ""},
		{"nil map, method", "{{.Method}}", nil, "", "can't evaluate field Method"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			err := Must(New(tc.name).Parse(tc.text)).ExecuteWithMethods(&b, methodMapData{"x"}, tc.methods)
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
	var b strings.Builder
	root := Must(New("root").Parse(`{{define "x"}}{{.Mapped}}{{end}}`))
	if err := root.Lookup("x").ExecuteWithMethods(&b, methodMapData{"<f>"}, methods); err != nil || b.String() != "&lt;&lt;f&gt;&gt;" {
		t.Fatalf("Lookup: got %q, %v; want the escaped result", b.String(), err)
	}
}
