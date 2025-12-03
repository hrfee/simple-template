# simple-template
[![Go Reference](https://pkg.go.dev/badge/github.com/hrfee/simple-template.svg)](https://pkg.go.dev/github.com/hrfee/simple-template)

simple templater for templates written by an end user. see godoc for more info.

## todo

- [ ] Implement else/else if
- [ ] Duplicate implementation in Typescript

## rough perf comparison

just for fun, not scientific. by suffix:
* "": This library
* "Old": Old version (-tags oldimpl)
* "BuiltinOnDemand": text/template with equivalent logic, template compile time included (-tags comparebuiltin)
* "BuiltinPrecompiled": text/template with equivalent logic, template compile time excluded (-tags comparebuiltin)
```
goos: linux
goarch: amd64
pkg: github.com/hrfee/simple-template
cpu: Intel(R) Core(TM) Ultra 9 185H
BenchmarkBlankTemplate-22                         	4023792	      287.4 ns/op
BenchmarkBlankTemplateOld-22                       	8906239	      135.1 ns/op
BenchmarkBlankTemplateBuiltinOnDemand-22          	 870487	     1175   ns/op
BenchmarkBlankTemplateBuiltinPrecompiled-22            10867418       107.5 ns/op
BenchmarkConditionalFalse-22                      	2090978	      571.5 ns/op
BenchmarkConditionalFalseOld-22                   	3605944	      335.6 ns/op
BenchmarkConditionalFalseBuiltinOnDemand-22       	 204522	     5780   ns/op
BenchmarkConditionalFalseBuiltinPrecompiled-22    	 797072	     1323   ns/op
BenchmarkConditionalTrue-22                       	1627858	      739.9 ns/op
BenchmarkConditionalTrueOld-22                    	1741509	      691.6 ns/op
BenchmarkConditionalTrueBuiltinOnDemand-22        	 195186	     5929   ns/op
BenchmarkConditionalTrueBuiltinPrecompiled-22     	1176372	     1017   ns/op
PASS
ok  	github.com/hrfee/simple-template	13.965s
```
