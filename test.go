package simpletemplate

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type templatingFunc func(in string, vars, conds []string, vals map[string]any) (string, error)
type newTemplatingFunc func(in string, vals map[string]any) (string, error)

func benchmarkBlankTemplate(b *testing.B, templateFunc templatingFunc) {
	in := `Success, user! Your account has been created. Log in at myAccountURL with your username to get started.`
	var vars, conds []string
	vals := map[string]any{}
	for b.Loop() {
		templateFunc(in, vars, conds, vals)
	}
}

func benchmarkConditional(isTrue bool, b *testing.B, templateFunc templatingFunc) {
	in := `Success, {username}! Your account has been created. {if myCondition}Log in at {myAccountURL} with username {username} to get started.{endif}`
	vars := []string{"username", "myAccountURL", "myCondition"}
	conds := vars
	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
		"myCondition":  isTrue,
	}
	for b.Loop() {
		templateFunc(in, vars, conds, vals)
	}
}

func benchmarkConditionalTrue(b *testing.B, templateFunc templatingFunc) {
	benchmarkConditional(true, b, templateFunc)
}
func benchmarkConditionalFalse(b *testing.B, templateFunc templatingFunc) {
	benchmarkConditional(false, b, templateFunc)
}

// In == Out when nothing is meant to be templated.
func testBlankTemplate(t *testing.T, templateFunc templatingFunc) {
	in := `Success, user! Your account has been created. Log in at myAccountURL with your username to get started.`

	out, err := templateFunc(in, []string{}, []string{}, map[string]any{})

	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	if out != in {
		t.Fatalf(`returned string doesn't match input: "%+v" != "%+v"`, out, in)
	}
}

func testConditional(isTrue bool, t *testing.T, templateFunc templatingFunc) {
	in := `Success, {username}! Your account has been created. {if myCondition}Log in at {myAccountURL} with username {username} to get started.{endif}`

	vars := []string{"username", "myAccountURL", "myCondition"}
	conds := vars
	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
		"myCondition":  isTrue,
	}

	out, err := templateFunc(in, vars, conds, vals)

	target := ""
	if isTrue {
		target = `Success, {username}! Your account has been created. Log in at {myAccountURL} with username {username} to get started.`
	} else {
		target = `Success, {username}! Your account has been created. `
	}

	target = strings.ReplaceAll(target, "{username}", vals["username"].(string))
	target = strings.ReplaceAll(target, "{myAccountURL}", vals["myAccountURL"].(string))

	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	if out != target {
		t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, target)
	}
}

func testConditionalTrue(t *testing.T, templateFunc templatingFunc) {
	testConditional(true, t, templateFunc)
}

func testConditionalFalse(t *testing.T, templateFunc templatingFunc) {
	testConditional(false, t, templateFunc)
}

// Template mistakenly double-braced values, but return a warning.
func testTemplateDoubleBraceGracefulHandling(t *testing.T, templateFunc templatingFunc) {
	in := `Success, {{username}}! Your account has been created. Log in at {myAccountURL} with username {username} to get started.`

	vars := []string{"username", "myAccountURL"}
	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
	}

	target := strings.ReplaceAll(in, "{{username}}", vals["username"].(string))
	target = strings.ReplaceAll(target, "{username}", vals["username"].(string))
	target = strings.ReplaceAll(target, "{myAccountURL}", vals["myAccountURL"].(string))

	out, err := templateFunc(in, vars, []string{}, vals)

	if err == nil {
		t.Fatal("no error when given double-braced variable")
	}

	if out != target {
		t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, target)
	}
}

func testVarAtAnyPosition(t *testing.T, templateFunc templatingFunc) {
	in := `Success, user! Your account has been created. Log in at myAccountURL with your username to get started.`
	vars := []string{"username", "myAccountURL"}
	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
	}

	for i := range in {
		newIn := in[0:i] + "{" + vars[0] + "}" + in[i:]

		target := strings.ReplaceAll(newIn, "{"+vars[0]+"}", vals["username"].(string))

		out, err := templateFunc(newIn, vars, []string{}, vals)

		if err != nil {
			t.Fatalf("error: %+v", err)
		}

		if out != target {
			t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v, from "%+v""`, out, target, newIn)
		}
	}
}

// In previous version, a lone { would be left alone but a warning would be returned.
// In new version, that's harder to implement so we'll return an error.
func testIncompleteBlock(t *testing.T, templateFunc templatingFunc) {
	in := `Success, user! Your account has been created. Log in at myAccountURL with your username to get started.`
	for i := range in {
		newIn := in[0:i] + "{" + in[i:]

		_, err := templateFunc(newIn, []string{"a"}, []string{"a"}, map[string]any{"a": "a"})

		// if out != newIn {
		// 	t.Fatalf(`returned string for position %d/%d doesn't match desired output: "%+v" != "%+v", err=%v`, i+1, len(newIn), out, newIn, err)
		// }
		if err == nil {
			t.Fatalf("no error when given incomplete block with brace at position %d/%d", i+1, len(newIn))
		}

	}
}

func testNegation(t *testing.T, templateFunc templatingFunc) {
	in := `Success, {username}! Your account has been created. {if !myCondition}Log in at {myAccountURL} with username {username} to get started.{endif}`

	vars := []string{"username", "myAccountURL", "myCondition"}
	conds := vars
	vals := map[string]any{
		"username":     "TemplateUsername",
		"myAccountURL": "TemplateURL",
	}

	f := func(isTrue bool, t *testing.T) {
		out, err := templateFunc(in, vars, conds, vals)

		target := ""
		if isTrue {
			target = `Success, {username}! Your account has been created. Log in at {myAccountURL} with username {username} to get started.`
		} else {
			target = `Success, {username}! Your account has been created. `
		}

		target = strings.ReplaceAll(target, "{username}", vals["username"].(string))
		target = strings.ReplaceAll(target, "{myAccountURL}", vals["myAccountURL"].(string))

		if err != nil {
			t.Fatalf("error: %+v", err)
		}

		if out != target {
			t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, target)
		}
	}

	t.Run("unassigned,true", func(t *testing.T) {
		f(true, t)
	})
	vals["myCondition"] = ""
	t.Run("blank,true", func(t *testing.T) {
		f(true, t)
	})
	vals["myCondition"] = "nonEmptyValue"
	t.Run("false", func(t *testing.T) {
		f(false, t)
	})
}

func testNestedIf(t *testing.T, templateFunc templatingFunc) {
	in := `{if varA}a{if varB}b{endif}{endif}`
	cases := []struct {
		name   string
		a, b   bool
		target string
	}{
		{"ff", false, false, ""},
		{"ft", false, true, ""},
		{"tf", true, false, "a"},
		{"tt", true, true, "ab"},
	}
	vars := []string{"varA", "varB"}
	conds := vars
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			vals := map[string]any{
				"varA": testCase.a,
				"varB": testCase.b,
			}
			out, err := templateFunc(in, vars, conds, vals)
			if err != nil {
				t.Fatalf("error: %+v", err)
			}
			if out != testCase.target {
				t.Fatalf(`returned string doesn't match desired output: "%+v" != "%+v"`, out, testCase.target)
			}
		})
	}
}

func testAdvancedConditional(t *testing.T, templateFunc newTemplatingFunc) {
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
			testAdvancedConditionalSub(testCase.isTrue, testCase.a, testCase.b, testCase.comp, t, templateFunc)
		})

	}
}
func testAdvancedConditionalSub(isTrue bool, valA, valB any, comp string, t *testing.T, templateFunc newTemplatingFunc) {
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
			out, err := templateFunc(testCase.input, vals)

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

func testSingleEqualsWarning(t *testing.T, templateFunc newTemplatingFunc) {
	in := `{if "myString" = "myString"}true{endif}`
	out, err := templateFunc(in, nil)

	var singleEquals SingleEqualsError
	if err == nil || !errors.As(err, &singleEquals) {
		t.Fatal("no error when given single equals in if block")
	}

	target := "true"
	if out != target {
		t.Fatalf(`template if not evaluated correctly when using a single equals: "%+v" != "%+v"`, out, target)
	}
}

func testIfElse(t *testing.T, templateFunc newTemplatingFunc) {
	t.Run("true", func(t *testing.T) { testIfElseSub(true, t, templateFunc) })
	t.Run("false", func(t *testing.T) { testIfElseSub(false, t, templateFunc) })
}

func testIfElseSub(isTrue bool, t *testing.T, templateFunc newTemplatingFunc) {
	in := `{if opA}a{else}b{endif}`
	out, err := templateFunc(in, map[string]any{
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

func testIfElseIf(t *testing.T, templateFunc newTemplatingFunc) {
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

			out, err := templateFunc(in, map[string]any{
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

func testAdvancedIfElseIf(t *testing.T, templateFunc newTemplatingFunc) {
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

			out, err := templateFunc(in, map[string]any{
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

func testIfElseIfElse(t *testing.T, templateFunc newTemplatingFunc) {
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

			out, err := templateFunc(in, map[string]any{
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
