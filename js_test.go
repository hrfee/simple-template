//go:build testjs

package simpletemplate

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"testing"
)

// So that tests can be run easily
func templateWrapperJS(in string, _, _ []string, vals map[string]any) (string, error) {
	return TemplateJS(in, vals)
}

func TemplateJS(in string, vals map[string]any) (string, error) {
	mapJSON, err := json.Marshal(vals)
	if err != nil {
		return "", err
	}
	var out, errBytes bytes.Buffer
	args := []string{"test_wrapper.js", in, string(mapJSON)}
	cmd := exec.Command("node", args...)
	cmd.Stderr = &errBytes
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	switch errBytes.String() {
	case "DoubleBraceError":
		err = DoubleBraceError{}
	case "SingleEqualsError":
		err = SingleEqualsError{}
	case "ExpectedTypeError":
		err = ExpectedTypeError{}
	case "ExpectedError":
		err = ExpectedError{}
	case "nil":
		err = nil
	default:
		err = errors.New(errBytes.String())
	}
	return out.String(), err
}

func TestBlankTemplateJS(t *testing.T)    { testBlankTemplate(t, templateWrapperJS) }
func TestConditionalTrueJS(t *testing.T)  { testConditionalTrue(t, templateWrapperJS) }
func TestConditionalFalseJS(t *testing.T) { testConditionalFalse(t, templateWrapperJS) }
func TestTemplateDoubleBraceGracefulHandlingJS(t *testing.T) {
	testTemplateDoubleBraceGracefulHandling(t, templateWrapperJS)
}
func TestVarAtAnyPositionJS(t *testing.T) { testVarAtAnyPosition(t, templateWrapperJS) }
func TestIncompleteBlockJS(t *testing.T)  { testIncompleteBlock(t, templateWrapperJS) }

func TestNegationJS(t *testing.T) { testNegation(t, templateWrapperJS) }

func TestAdvancedConditionalJS(t *testing.T) { testAdvancedConditional(t, TemplateJS) }
func TestSingleEqualsWarningJS(t *testing.T) { testSingleEqualsWarning(t, TemplateJS) }
func TestNestedIfJS(t *testing.T)            { testNestedIf(t, templateWrapperJS) }
func TestIfElseJS(t *testing.T)              { testIfElse(t, TemplateJS) }
func TestIfElseIfJS(t *testing.T)            { testIfElseIf(t, TemplateJS) }
func TestAdvancedIfElseIfJS(t *testing.T)    { testAdvancedIfElseIf(t, TemplateJS) }
func TestIfElseIfElseJS(t *testing.T)        { testIfElseIfElse(t, TemplateJS) }
