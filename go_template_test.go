//go:build comparebuiltin

package simpletemplate

import (
	"strings"
	"testing"
	"text/template"
)

func BenchmarkBlankTemplateBuiltinOnDemand(b *testing.B) {
	in := `Success, user! Your account has been created. Log in at myAccountURL with your username to get started.`
	vals := map[string]any{}
	for b.Loop() {
		var out strings.Builder
		out.Grow(len(in))
		tmpl, _ := template.New("test").Parse(in)
		tmpl.Execute(&out, vals)
	}
}
func benchmarkConditionalBuiltinOnDemand(isTrue bool, b *testing.B) {
	in := `Success, {{ .username }}! Your account has been created. {{ if not eq .myCondition "" }}Log in at {{ .myAccountURL }} with username {{ .username }} to get started.{{ end }}`
	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
		"myCondition":  isTrue,
	}
	for b.Loop() {
		var out strings.Builder
		out.Grow(len(in))
		tmpl, _ := template.New("test").Parse(in)
		tmpl.Execute(&out, vals)
	}
}
func BenchmarkConditionalTrueBuiltinOnDemand(b *testing.B) {
	benchmarkConditionalBuiltinOnDemand(true, b)
}
func BenchmarkConditionalFalseBuiltinOnDemand(b *testing.B) {
	benchmarkConditionalBuiltinOnDemand(false, b)
}

func BenchmarkBlankTemplateBuiltinPrecompiled(b *testing.B) {
	in := `Success, user! Your account has been created. Log in at myAccountURL with your username to get started.`
	vals := map[string]any{}
	tmpl, _ := template.New("test").Parse(in)
	for b.Loop() {
		var out strings.Builder
		out.Grow(len(in))
		tmpl.Execute(&out, vals)
	}
}
func benchmarkConditionalBuiltinPrecompiled(isTrue bool, b *testing.B) {
	in := `Success, {{ .username }}! Your account has been created. {{ if not (eq .myCondition true) }}Log in at {{ .myAccountURL }} with username {{ .username }} to get started.{{ end }}`
	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
		"myCondition":  isTrue,
	}
	tmpl, _ := template.New("test").Parse(in)
	for b.Loop() {
		var out strings.Builder
		out.Grow(len(in))
		tmpl.Execute(&out, vals)
	}
}
func BenchmarkConditionalTrueBuiltinPrecompiled(b *testing.B) {
	benchmarkConditionalBuiltinPrecompiled(true, b)
}
func BenchmarkConditionalFalseBuiltinPrecompiled(b *testing.B) {
	benchmarkConditionalBuiltinPrecompiled(false, b)
}
