# Jet Template Engine for Go [![Build Status](https://travis-ci.org/CloudyKit/jet.svg?branch=master)](https://travis-ci.org/CloudyKit/jet) [![Build status](https://ci.appveyor.com/api/projects/status/5g4whw3c6518vvku?svg=true)](https://ci.appveyor.com/project/CloudyKit/jet)

[![Join the chat at https://gitter.im/CloudyKit/jet](https://badges.gitter.im/CloudyKit/jet.svg)](https://gitter.im/CloudyKit/jet)

Jet is a template engine developed to be easy to use, powerful, dynamic, yet secure and very fast.

* supports template inheritance with `extends`, `import` and `include` statements
* descriptive error messages with filename and line number
* auto-escape
* simple C-like expressions
* very fast execution – Jet can execute templates faster than some pre-compiled template engines
* very light in terms of allocations and memory footprint
* simple and familiar syntax
* easy to use

You can find the documentation in the [wiki](https://github.com/CloudyKit/jet/wiki).

#### Upgrade to v2

The last release of v1 was v1.2 which is available at https://github.com/CloudyKit/jet/releases/tag/v1.2 and the tag v1.2.

To upgrade to v2 a few updates to your templates are necessary – these are explained in the [upgrade guide](https://github.com/CloudyKit/jet/wiki/Upgrade-to-v2).

#### IntelliJ Plugin

If you use IntelliJ there is a plugin available at https://github.com/jhsx/GoJetPlugin.
There is also a very good Go plugin for IntelliJ – see https://github.com/go-lang-plugin-org/go-lang-idea-plugin.
GoJetPlugin + Go-lang-idea-plugin = happiness!

### Examples

You can find examples in the [wiki](https://github.com/CloudyKit/jet/wiki/Jet-template-syntax).

### Running the example application

An example application is available in the repository. Use `go get -u github.com/CloudyKit/jet` or clone the repository into `$GOPATH/github.com/CloudyKit/jet`, then do:
```
  $ cd examples/todos; go run main.go
```

#### Faster than some pre-compiled template engines

The benchmark consists of a range over a slice of data printing the values, the benchmark is based on https://github.com/SlinSo/goTemplateBenchmark, Jet performs better than all template engines without pre-compilation,
and performs better than gorazor, Ftmpl and Egon, all of which are pre-compiled to Go.

###### Benchmarks

go 1.6.2
```
PASS
BenchmarkEgonSlinso-4      	 2000000	       989 ns/op	     517 B/op	       0 allocs/op
BenchmarkQuicktemplate-4   	 1000000	      1004 ns/op	     999 B/op	       0 allocs/op
BenchmarkEgo-4             	 1000000	      2137 ns/op	     603 B/op	       8 allocs/op

BenchmarkJet-4             	  500000	      2798 ns/op	     536 B/op	       0 allocs/op
BenchmarkJetHTML-4         	  500000	      2822 ns/op	     536 B/op	       0 allocs/op

BenchmarkGorazor-4         	  500000	      3028 ns/op	     613 B/op	      11 allocs/op
BenchmarkFtmpl-4           	  500000	      3192 ns/op	    1152 B/op	      12 allocs/op
BenchmarkEgon-4            	  300000	      4673 ns/op	    1172 B/op	      22 allocs/op
BenchmarkKasia-4           	  200000	      6902 ns/op	    1789 B/op	      26 allocs/op
BenchmarkSoy-4             	  200000	      7144 ns/op	    1684 B/op	      26 allocs/op
BenchmarkMustache-4        	  200000	      8213 ns/op	    1568 B/op	      28 allocs/op
BenchmarkPongo2-4          	  200000	      9989 ns/op	    2949 B/op	      46 allocs/op
BenchmarkGolang-4          	  100000	     16284 ns/op	    2039 B/op	      38 allocs/op
BenchmarkAmber-4           	  100000	     17208 ns/op	    2050 B/op	      39 allocs/op
BenchmarkHandlebars-4      	   50000	     29864 ns/op	    4258 B/op	      90 allocs/op
BenchmarkAce-4             	   30000	     40771 ns/op	    5710 B/op	      77 allocs/op
BenchmarkDamsel-4          	   20000	     95947 ns/op	   11160 B/op	     165 allocs/op
ok  	github.com/SlinSo/goTemplateBenchmark	34.384s
```

go tip
```
BenchmarkQuicktemplate-4      	 2000000	       916 ns/op	     999 B/op	       0 allocs/op
BenchmarkEgonSlinso-4         	 2000000	      1074 ns/op	     517 B/op	       0 allocs/op
BenchmarkEgo-4                	 1000000	      1822 ns/op	     603 B/op	       8 allocs/op

BenchmarkJetHTML-4            	  500000	      2627 ns/op	     536 B/op	       0 allocs/op
BenchmarkJet-4                	  500000	      2652 ns/op	     536 B/op	       0 allocs/op

BenchmarkFtmpl-4              	  500000	      2700 ns/op	    1152 B/op	      12 allocs/op
BenchmarkGorazor-4            	  500000	      2858 ns/op	     613 B/op	      11 allocs/op
BenchmarkEgon-4               	  500000	      4023 ns/op	     827 B/op	      22 allocs/op
BenchmarkSoy-4                	  300000	      5590 ns/op	    1784 B/op	      26 allocs/op
BenchmarkKasia-4              	  200000	      6487 ns/op	    1789 B/op	      26 allocs/op
BenchmarkMustache-4           	  200000	      6515 ns/op	    1568 B/op	      28 allocs/op
BenchmarkPongo2-4             	  200000	      7602 ns/op	    2949 B/op	      46 allocs/op
BenchmarkAmber-4              	  100000	     13942 ns/op	    2050 B/op	      39 allocs/op
BenchmarkGolang-4             	  100000	     16945 ns/op	    2039 B/op	      38 allocs/op
BenchmarkHandlebars-4         	  100000	     20152 ns/op	    4258 B/op	      90 allocs/op
BenchmarkAce-4                	   50000	     33091 ns/op	    5509 B/op	      77 allocs/op
BenchmarkDamsel-4             	   20000	     86340 ns/op	   11159 B/op	     165 allocs/op
PASS
ok  	github.com/SlinSo/goTemplateBenchmark	36.200s
```

#### Contributing

All contributions are welcome – if you find a bug please report it.

#### Thanks

- @golang developers for the awesome language and the standard library
- @SlinSo for the benchmarks that I used as a base to show the results above
