package neffos

import (
	"log"
	"os"
	"reflect"
)

var debugPrinter interface{}

// EnableDebug enables debug and optionally
// sets a custom printer to print out debug messages.
// The "printer" can be any compatible printer such as
// the standard `log.Logger` or a custom one like the `kataras/golog`.
//
// A "printer" is compatible when it contains AT LEAST ONE of the following methods:
// Debugf(string, ...interface{}) or
// Logf(string, ...interface{}) or
// Printf(string, ...interface{})
//
// If EnableDebug is called but the "printer" value is nil
// then neffos will print debug messages through a new log.Logger prefixed with "| neffos |".
//
// Note that neffos, currently, uses debug mode only on the build state of the events.
// Therefore enabling the debugger has zero performance cost on up-and-running servers and clients.
//
// There is no way to disable the debug mode on serve-time.
func EnableDebug(printer interface{}) {
	if debugEnabled() {
		Debugf("debug mode is already set")
		return
	}

	if _, boolean := printer.(bool); boolean {
		// if for some reason by accident EnableDebug(true) instead of a printer value.
		printer = nil
	}
	if printer == nil {
		logger := log.New(os.Stderr, "| neffos | ", 0)
		printer = logger
		logger.Println("debug mode is set")
	}

	debugPrinter = printer
}

type (
	debugfer interface {
		Debugf(string, ...interface{})
	}
	logfer interface {
		Logf(string, ...interface{})
	}
	printfer interface {
		Printf(string, ...interface{})
	}
)

func debugEnabled() bool {
	return debugPrinter != nil
}

// Debugf prints debug messages to the printer defined on `EnableDebug`.
// Runs only on debug mode.
func Debugf(format string, args ...interface{}) {
	if !debugEnabled() {
		return
	}

	if len(args) == 1 {
		// handles:
		// Debugf("format", func() dargs {
		//    time-consumed action that should run only on debug.
		// })
		if onDebugWithArgs, ok := args[0].(func() dargs); ok {
			args = onDebugWithArgs()
		}
	}

	switch printer := debugPrinter.(type) {
	case debugfer:
		printer.Debugf(format, args...)
	case logfer:
		printer.Logf(format, args...)
	case printfer:
		printer.Printf(format, args...)
	default:
		panic("unsported debug printer")
	}
}

type dargs []interface{}

// DebugEach prints debug messages for each of "mapOrSlice" elements
// to the printer defined on `EnableDebug`.
// Runs only on debug mode.
// Usage:
// DebugEach(staticFields, func(idx int, f reflect.Value) {
// 	fval := f.Interface()
// 	Debugf("field [%s.%s] will be automatically re-filled with [%T(%s)]", typ.Name(), typ.Field(idx).Name, fval, fval)
// })
func DebugEach(mapOrSlice interface{}, onDebugVisitor interface{}) {
	if !debugEnabled() || onDebugVisitor == nil {
		return
	}

	visitor := reflect.ValueOf(onDebugVisitor)

	visitorTyp := visitor.Type()
	if visitorTyp.Kind() != reflect.Func {
		return
	}

	userNumIn := visitorTyp.NumIn()

	v := reflect.ValueOf(mapOrSlice)

	switch v.Kind() {
	case reflect.Map:
		for ranger := v.MapRange(); ranger.Next(); {
			in := make([]reflect.Value, userNumIn)
			in[0] = ranger.Key()

			if userNumIn > 1 {
				// assume both key and value are expected.
				in[1] = ranger.Value()
			}

			// note that we don't make any further checks here, it's only for internal
			// use and I want to panic in my tests if I didn't expect the correct values.
			visitor.Call(in)
		}
	case reflect.Slice:
		// TODO: whenever I need this.
	}
}
