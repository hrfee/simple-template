package simpletemplate

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// So that tests can be run easily comparing the old and new implementation
func templateWrapper(in string, _, _ []string, vals map[string]any) (string, error) {
	return Template(in, vals)
}

func BenchmarkBlankTemplate(b *testing.B)    { benchmarkBlankTemplate(b, templateWrapper) }
func BenchmarkConditionalTrue(b *testing.B)  { benchmarkConditionalTrue(b, templateWrapper) }
func BenchmarkConditionalFalse(b *testing.B) { benchmarkConditionalFalse(b, templateWrapper) }

func TestBlankTemplate(t *testing.T)    { testBlankTemplate(t, templateWrapper) }
func TestConditionalTrue(t *testing.T)  { testConditionalTrue(t, templateWrapper) }
func TestConditionalFalse(t *testing.T) { testConditionalFalse(t, templateWrapper) }
func TestTemplateDoubleBraceGracefulHandling(t *testing.T) {
	testTemplateDoubleBraceGracefulHandling(t, templateWrapper)
}
func TestVarAtAnyPosition(t *testing.T) { testVarAtAnyPosition(t, templateWrapper) }
func TestIncompleteBlock(t *testing.T)  { testIncompleteBlock(t, templateWrapper) }

func TestNegation(t *testing.T) { testNegation(t, templateWrapper) }

func TestAdvancedConditional(t *testing.T) {
	cases := []struct {
		name   string
		a, b   any
		comp   string
		isTrue bool
	}{
		{"a==aT", "a string", "a string", "==", true},
		{"a==bF", "a string", "b string", "==", false},
		{"a!=bT", "a string", "b string", "!=", true},
		{"a!=aF", "a string", "a string", "!=", false},
	}

	for _, testCase := range cases {
		if testCase.comp == "==" {
			testCase.comp = "="
			testCase.name = strings.Replace(testCase.name, "==", "=", 1)
			cases = append(cases, testCase)
		}
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testAdvancedConditional(testCase.isTrue, testCase.a, testCase.b, testCase.comp, t)
		})

	}
}

// The following tests are for features not present in the old implementation.
func testAdvancedConditional(isTrue bool, valA, valB any, comp string, t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"var,var", fmt.Sprintf(`Success, {username}! Your account has been created. {if opA %s opB}Log in at {myAccountURL} with username {username} to get started.{endif}`, comp)},
		{"literal,var", fmt.Sprintf(`Success, {username}! Your account has been created. {if "%s" %s opB}Log in at {myAccountURL} with username {username} to get started.{endif}`, valA, comp)},
		{"var,literal", fmt.Sprintf(`Success, {username}! Your account has been created. {if opA %s "%s"}Log in at {myAccountURL} with username {username} to get started.{endif}`, comp, valB)},
		{"literal,literal", fmt.Sprintf(`Success, {username}! Your account has been created. {if '%s' %s "%s"}Log in at {myAccountURL} with username {username} to get started.{endif}`, valA, comp, valB)},
	}

	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
		"opA":          valA,
		"opB":          valB,
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			out, err := Template(testCase.input, vals)

			var target string
			if isTrue {
				target = `Success, {username}! Your account has been created. Log in at {myAccountURL} with username {username} to get started.`
			} else {
				target = `Success, {username}! Your account has been created. `
			}

			target = strings.ReplaceAll(target, "{username}", vals["username"].(string))
			target = strings.ReplaceAll(target, "{myAccountURL}", vals["myAccountURL"].(string))

			if err != nil {
				var singleEquals SingleEqualsError
				if comp != "=" || !errors.As(err, &singleEquals) {
					t.Fatalf("error: %+v", err)
				}
			}

			if out != target {
				t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, target)
			}
		})
	}
}

func TestSingleEqualsWarning(t *testing.T) {
	in := `{if "myString" = "myString"}true{endif}`
	out, err := Template(in, nil)

	var singleEquals SingleEqualsError
	if err == nil || !errors.As(err, &singleEquals) {
		t.Fatal("no error when given single equals in if block")
	}

	target := "true"
	if out != target {
		t.Fatalf(`template if not evaluated correctly when using a single equals: "%+v" != "%+v"`, out, target)
	}
}

func TestNestedIf(t *testing.T) { testNestedIf(t, templateWrapper) }

func TestIfElse(t *testing.T) {
	t.Run("true", func(t *testing.T) { testIfElse(true, t) })
	t.Run("false", func(t *testing.T) { testIfElse(false, t) })
}

func testIfElse(isTrue bool, t *testing.T) {
	in := `{if opA}a{else}b{endif}`
	out, err := Template(in, map[string]any{
		"opA": isTrue,
	})

	var target string
	if isTrue {
		target = "a"
	} else {
		target = "b"
	}

	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	if out != target {
		t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, target)
	}
}

func TestIfElseIf(t *testing.T) {
	in := `{if opA}a{else if opB}b{endif}`
	cases := []struct {
		name   string
		a, b   bool
		target string
	}{
		{"ff", false, false, ""},
		{"ft", false, true, "b"},
		{"tf", true, false, "a"},
		{"tt", true, true, "a"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {

			out, err := Template(in, map[string]any{
				"opA": testCase.a,
				"opB": testCase.b,
			})

			if err != nil {
				t.Fatalf("error: %+v", err)
			}

			if out != testCase.target {
				t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, testCase.target)
			}
		})
	}
}

func TestAdvancedIfElseIf(t *testing.T) {
	in := `{if opA == "a"}a{else if opB != "b"}b{endif}`
	cases := []struct {
		a, b   string
		target string
	}{
		{"a", "a", "a"},
		{"a", "b", "a"},
		{"b", "a", "b"},
		{"b", "b", ""},
	}
	for _, testCase := range cases {
		t.Run(testCase.a+testCase.b, func(t *testing.T) {

			out, err := Template(in, map[string]any{
				"opA": testCase.a,
				"opB": testCase.b,
			})

			if err != nil {
				t.Fatalf("error: %+v", err)
			}

			if out != testCase.target {
				t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, testCase.target)
			}
		})
	}
}

func TestIfElseIfElse(t *testing.T) {
	in := `{if opA}a{else if opB}b{else}c{endif}`
	cases := []struct {
		name   string
		a, b   bool
		target string
	}{
		{"ff", false, false, "c"},
		{"ft", false, true, "b"},
		{"tf", true, false, "a"},
		{"tt", true, true, "a"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {

			out, err := Template(in, map[string]any{
				"opA": testCase.a,
				"opB": testCase.b,
			})

			if err != nil {
				t.Fatalf("error: %+v", err)
			}

			if out != testCase.target {
				t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, testCase.target)
			}
		})
	}
}
