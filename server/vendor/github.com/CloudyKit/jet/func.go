// Copyright 2016 José Santos <henrique_1609@me.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jet

import (
	"fmt"
	"reflect"
)

// Arguments holds the arguments passed to jet.Func.
type Arguments struct {
	runtime *Runtime
	argExpr []Expression
	argVal  []reflect.Value
}

// Get gets an argument by index.
func (a *Arguments) Get(argumentIndex int) reflect.Value {
	if argumentIndex < len(a.argVal) {
		return a.argVal[argumentIndex]
	}
	if argumentIndex < len(a.argVal)+len(a.argExpr) {
		return a.runtime.evalPrimaryExpressionGroup(a.argExpr[argumentIndex-len(a.argVal)])
	}
	return reflect.Value{}
}

// Panicf panics with formatted error message.
func (a *Arguments) Panicf(format string, v ...interface{}) {
	panic(fmt.Errorf(format, v...))
}

// RequireNumOfArguments panics if the number of arguments is not in the range specified by min and max.
// In case there is no minimum pass -1, in case there is no maximum pass -1 respectively.
func (a *Arguments) RequireNumOfArguments(funcname string, min, max int) {
	num := len(a.argExpr) + len(a.argVal)
	if min >= 0 && num < min {
		a.Panicf("unexpected number of arguments in a call to %s", funcname)
	} else if max >= 0 && num > max {
		a.Panicf("unexpected number of arguments in a call to %s", funcname)
	}
}

// NumOfArguments returns the number of arguments
func (a *Arguments) NumOfArguments() int {
	return len(a.argExpr) + len(a.argVal)
}

// Runtime get the Runtime context
func (a *Arguments) Runtime() *Runtime {
	return a.runtime
}

// Func function implementing this type is called directly, which is faster than calling through reflect.
// If a function is being called many times in the execution of a template, you may consider implementing
// a wrapper for that function implementing a Func.
type Func func(Arguments) reflect.Value
