package simpletemplate_test

import (
	"fmt"

	simpletemplate "github.com/hrfee/simple-template"
)

func Example() {
	in := `
	Here is some plain text. The value of variable varA is {varA}.
	{if varB == "true"}varB is set to true{else}varB is not set to true, it's set to {varB}.{endif}
	{if varC != "1"}varC is not 1{else if varD == "1"}varC and varD are set to 1.{endif} 
	{if varE}varE has some non-empty value set.{endif}
	{if !varF}varF is unset or a zero-value.{endif}
	{if !varG}varG is unset or a zero-value.{endif}
	`

	out, err := simpletemplate.Template(in, map[string]any{
		"varA": "aValue",
		"varB": "true",
		"varC": 1,
		"varD": "1",
		"varE": "anotherValue",
		"varF": 0,
	})
	fmt.Printf("out: \"%s\"\n", out)
	fmt.Printf("err: %v\n", err)
	// out: "
	// 	Here is some plain text. The value of variable varA is aValue.
	// 	varB is set to true
	// 	varC is not 1
	// 	varE has some non-empty value set.
	// 	varF is unset or a zero-value.
	// 	varG is unset or a zero-value.
	// 	"
	// err: <nil>
}
