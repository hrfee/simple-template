package simpletemplate_test

import (
	"fmt"

	simpletemplate "github.com/hrfee/simple-template"
)

func Example() {
	in := `
	Here is some plain text. The value of variable varA is {varA}.
	{if varB == "true"}varB is set to true{endif}{if varB != "true"}varB is not set to true, it's set to {varB}.{endif}
	{if varC}varC has some non-empty value set.{endif}
	{if !varD}varD is unset or a zero-value.{endif}`

	out, err := simpletemplate.Template(in, map[string]any{
		"varA": "aValue",
		"varB": "true",
		"varC": "anotherValue",
		"varD": 0,
	})
	fmt.Printf("out: \"%s\"\n", out)
	fmt.Printf("err: %v\n", err)
	// Output:
	// out: "
	// 	Here is some plain text. The value of variable varA is aValue.
	// 	varB is set to true
	// 	varC has some non-empty value set.
	// 	varD is unset or a zero-value."
	// err: <nil>
}
