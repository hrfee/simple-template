//go:build oldimpl

package simpletemplate

import "testing"

func BenchmarkBlankTemplateOld(b *testing.B)    { benchmarkBlankTemplate(b, TemplateOld) }
func BenchmarkConditionalTrueOld(b *testing.B)  { benchmarkConditionalTrue(b, TemplateOld) }
func BenchmarkConditionalFalseOld(b *testing.B) { benchmarkConditionalFalse(b, TemplateOld) }

func TestBlankTemplateOld(t *testing.T)    { testBlankTemplate(t, TemplateOld) }
func TestConditionalTrueOld(t *testing.T)  { testConditionalTrue(t, TemplateOld) }
func TestConditionalFalseOld(t *testing.T) { testConditionalFalse(t, TemplateOld) }
func TestTemplateDoubleBraceGracefulHandlingOld(t *testing.T) {
	testTemplateDoubleBraceGracefulHandling(t, TemplateOld)
}
func TestVarAtAnyPositionOld(t *testing.T) { testVarAtAnyPosition(t, TemplateOld) }
func TestIncompleteBlockOld(t *testing.T)  { testIncompleteBlock(t, TemplateOld) }

func TestNegationOld(t *testing.T) { testNegation(t, TemplateOld) }

func TestNestedIfOld(t *testing.T) { testNestedIf(t, TemplateOld) }
