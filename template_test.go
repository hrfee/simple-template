package simpletemplate

import (
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

func TestAdvancedConditional(t *testing.T) { testAdvancedConditional(t, Template) }
func TestSingleEqualsWarning(t *testing.T) { testSingleEqualsWarning(t, Template) }
func TestNestedIf(t *testing.T)            { testNestedIf(t, templateWrapper) }
func TestIfElse(t *testing.T)              { testIfElse(t, Template) }
func TestIfElseIf(t *testing.T)            { testIfElseIf(t, Template) }
func TestAdvancedIfElseIf(t *testing.T)    { testAdvancedIfElseIf(t, Template) }
func TestIfElseIfElse(t *testing.T)        { testIfElseIfElse(t, Template) }
