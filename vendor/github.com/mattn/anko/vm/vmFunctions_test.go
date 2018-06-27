package vm

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mattn/anko/internal/testlib"
)

func TestReturns(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testlib.Test{
		{Script: `return 1++`, RunError: fmt.Errorf("invalid operation")},
		{Script: `return 1, 1++`, RunError: fmt.Errorf("invalid operation")},
		{Script: `return 1, 2, 1++`, RunError: fmt.Errorf("invalid operation")},

		{Script: `return`, RunOutput: nil},
		{Script: `return nil`, RunOutput: nil},
		{Script: `return true`, RunOutput: true},
		{Script: `return 1`, RunOutput: int64(1)},
		{Script: `return 1.1`, RunOutput: float64(1.1)},
		{Script: `return "a"`, RunOutput: "a"},

		{Script: `b()`, Input: map[string]interface{}{"b": func() {}}, RunOutput: nil},
		{Script: `b()`, Input: map[string]interface{}{"b": func() reflect.Value { return reflect.Value{} }}, RunOutput: reflect.Value{}},
		{Script: `b()`, Input: map[string]interface{}{"b": func() interface{} { return nil }}, RunOutput: nil},
		{Script: `b()`, Input: map[string]interface{}{"b": func() bool { return true }}, RunOutput: true},
		{Script: `b()`, Input: map[string]interface{}{"b": func() int32 { return int32(1) }}, RunOutput: int32(1)},
		{Script: `b()`, Input: map[string]interface{}{"b": func() int64 { return int64(1) }}, RunOutput: int64(1)},
		{Script: `b()`, Input: map[string]interface{}{"b": func() float32 { return float32(1.1) }}, RunOutput: float32(1.1)},
		{Script: `b()`, Input: map[string]interface{}{"b": func() float64 { return float64(1.1) }}, RunOutput: float64(1.1)},
		{Script: `b()`, Input: map[string]interface{}{"b": func() string { return "a" }}, RunOutput: "a"},

		{Script: `b(a)`, Input: map[string]interface{}{"a": reflect.Value{}, "b": func(c reflect.Value) reflect.Value { return c }}, RunOutput: reflect.Value{}, Output: map[string]interface{}{"a": reflect.Value{}}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": nil, "b": func(c interface{}) interface{} { return c }}, RunOutput: nil, Output: map[string]interface{}{"a": nil}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": true, "b": func(c bool) bool { return c }}, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": int32(1), "b": func(c int32) int32 { return c }}, RunOutput: int32(1), Output: map[string]interface{}{"a": int32(1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": int64(1), "b": func(c int64) int64 { return c }}, RunOutput: int64(1), Output: map[string]interface{}{"a": int64(1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": float32(1.1), "b": func(c float32) float32 { return c }}, RunOutput: float32(1.1), Output: map[string]interface{}{"a": float32(1.1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": float64(1.1), "b": func(c float64) float64 { return c }}, RunOutput: float64(1.1), Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": "a", "b": func(c string) string { return c }}, RunOutput: "a", Output: map[string]interface{}{"a": "a"}},

		{Script: `b(a)`, Input: map[string]interface{}{"a": "a", "b": func(c bool) bool { return c }}, RunError: fmt.Errorf("function wants argument type bool but received type string"), Output: map[string]interface{}{"a": "a"}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": int64(1), "b": func(c int32) int32 { return c }}, RunOutput: int32(1), Output: map[string]interface{}{"a": int64(1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": int32(1), "b": func(c int64) int64 { return c }}, RunOutput: int64(1), Output: map[string]interface{}{"a": int32(1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": float64(1.25), "b": func(c float32) float32 { return c }}, RunOutput: float32(1.25), Output: map[string]interface{}{"a": float64(1.25)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": float32(1.25), "b": func(c float64) float64 { return c }}, RunOutput: float64(1.25), Output: map[string]interface{}{"a": float32(1.25)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": true, "b": func(c string) string { return c }}, RunError: fmt.Errorf("function wants argument type string but received type bool"), Output: map[string]interface{}{"a": true}},

		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueBool, "b": func(c interface{}) interface{} { return c }}, RunOutput: testVarValueBool, Output: map[string]interface{}{"a": testVarValueBool}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueInt32, "b": func(c interface{}) interface{} { return c }}, RunOutput: testVarValueInt32, Output: map[string]interface{}{"a": testVarValueInt32}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueInt64, "b": func(c interface{}) interface{} { return c }}, RunOutput: testVarValueInt64, Output: map[string]interface{}{"a": testVarValueInt64}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueFloat32, "b": func(c interface{}) interface{} { return c }}, RunOutput: testVarValueFloat32, Output: map[string]interface{}{"a": testVarValueFloat32}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueFloat64, "b": func(c interface{}) interface{} { return c }}, RunOutput: testVarValueFloat64, Output: map[string]interface{}{"a": testVarValueFloat64}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueString, "b": func(c interface{}) interface{} { return c }}, RunOutput: testVarValueString, Output: map[string]interface{}{"a": testVarValueString}},

		{Script: `func aFunc() {}; aFunc()`, RunOutput: nil},
		{Script: `func aFunc() {return}; aFunc()`, RunOutput: nil},
		{Script: `func aFunc() {return}; a = aFunc()`, RunOutput: nil, Output: map[string]interface{}{"a": nil}},

		{Script: `func aFunc() {return nil}; aFunc()`, RunOutput: nil},
		{Script: `func aFunc() {return true}; aFunc()`, RunOutput: true},
		{Script: `func aFunc() {return 1}; aFunc()`, RunOutput: int64(1)},
		{Script: `func aFunc() {return 1.1}; aFunc()`, RunOutput: float64(1.1)},
		{Script: `func aFunc() {return "a"}; aFunc()`, RunOutput: "a"},

		{Script: `func aFunc() {return 1 + 2}; aFunc()`, RunOutput: int64(3)},
		{Script: `func aFunc() {return 1.25 + 2.25}; aFunc()`, RunOutput: float64(3.5)},
		{Script: `func aFunc() {return "a" + "b"}; aFunc()`, RunOutput: "ab"},

		{Script: `func aFunc() {return 1 + 2, 3 + 4}; aFunc()`, RunOutput: []interface{}{int64(3), int64(7)}},
		{Script: `func aFunc() {return 1.25 + 2.25, 3.25 + 4.25}; aFunc()`, RunOutput: []interface{}{float64(3.5), float64(7.5)}},
		{Script: `func aFunc() {return "a" + "b", "c" + "d"}; aFunc()`, RunOutput: []interface{}{"ab", "cd"}},

		{Script: `func aFunc() {return nil, nil}; aFunc()`, RunOutput: []interface{}{nil, nil}},
		{Script: `func aFunc() {return true, false}; aFunc()`, RunOutput: []interface{}{true, false}},
		{Script: `func aFunc() {return 1, 2}; aFunc()`, RunOutput: []interface{}{int64(1), int64(2)}},
		{Script: `func aFunc() {return 1.1, 2.2}; aFunc()`, RunOutput: []interface{}{float64(1.1), float64(2.2)}},
		{Script: `func aFunc() {return "a", "b"}; aFunc()`, RunOutput: []interface{}{"a", "b"}},

		{Script: `func aFunc() {return [nil]}; aFunc()`, RunOutput: []interface{}{nil}},
		{Script: `func aFunc() {return [nil, nil]}; aFunc()`, RunOutput: []interface{}{nil, nil}},
		{Script: `func aFunc() {return [nil, nil, nil]}; aFunc()`, RunOutput: []interface{}{nil, nil, nil}},
		{Script: `func aFunc() {return [nil, nil], [nil, nil]}; aFunc()`, RunOutput: []interface{}{[]interface{}{nil, nil}, []interface{}{nil, nil}}},

		{Script: `func aFunc() {return [true]}; aFunc()`, RunOutput: []interface{}{true}},
		{Script: `func aFunc() {return [true, false]}; aFunc()`, RunOutput: []interface{}{true, false}},
		{Script: `func aFunc() {return [true, false, true]}; aFunc()`, RunOutput: []interface{}{true, false, true}},
		{Script: `func aFunc() {return [true, false], [false, true]}; aFunc()`, RunOutput: []interface{}{[]interface{}{true, false}, []interface{}{false, true}}},

		{Script: `func aFunc() {return []}; aFunc()`, RunOutput: []interface{}{}},
		{Script: `func aFunc() {return [1]}; aFunc()`, RunOutput: []interface{}{int64(1)}},
		{Script: `func aFunc() {return [1, 2]}; aFunc()`, RunOutput: []interface{}{int64(1), int64(2)}},
		{Script: `func aFunc() {return [1, 2, 3]}; aFunc()`, RunOutput: []interface{}{int64(1), int64(2), int64(3)}},
		{Script: `func aFunc() {return [1, 2], [3, 4]}; aFunc()`, RunOutput: []interface{}{[]interface{}{int64(1), int64(2)}, []interface{}{int64(3), int64(4)}}},

		{Script: `func aFunc() {return [1.1]}; aFunc()`, RunOutput: []interface{}{float64(1.1)}},
		{Script: `func aFunc() {return [1.1, 2.2]}; aFunc()`, RunOutput: []interface{}{float64(1.1), float64(2.2)}},
		{Script: `func aFunc() {return [1.1, 2.2, 3.3]}; aFunc()`, RunOutput: []interface{}{float64(1.1), float64(2.2), float64(3.3)}},
		{Script: `func aFunc() {return [1.1, 2.2], [3.3, 4.4]}; aFunc()`, RunOutput: []interface{}{[]interface{}{float64(1.1), float64(2.2)}, []interface{}{float64(3.3), float64(4.4)}}},

		{Script: `func aFunc() {return ["a"]}; aFunc()`, RunOutput: []interface{}{"a"}},
		{Script: `func aFunc() {return ["a", "b"]}; aFunc()`, RunOutput: []interface{}{"a", "b"}},
		{Script: `func aFunc() {return ["a", "b", "c"]}; aFunc()`, RunOutput: []interface{}{"a", "b", "c"}},
		{Script: `func aFunc() {return ["a", "b"], ["c", "d"]}; aFunc()`, RunOutput: []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}}},

		{Script: `func aFunc() {return nil, nil}; aFunc()`, RunOutput: []interface{}{interface{}(nil), interface{}(nil)}},
		{Script: `func aFunc() {return true, false}; aFunc()`, RunOutput: []interface{}{true, false}},
		{Script: `func aFunc() {return 1, 2}; aFunc()`, RunOutput: []interface{}{int64(1), int64(2)}},
		{Script: `func aFunc() {return 1.1, 2.2}; aFunc()`, RunOutput: []interface{}{float64(1.1), float64(2.2)}},
		{Script: `func aFunc() {return "a", "b"}; aFunc()`, RunOutput: []interface{}{"a", "b"}},

		{Script: `func aFunc() {return a}; aFunc()`, Input: map[string]interface{}{"a": reflect.Value{}}, RunOutput: reflect.Value{}, Output: map[string]interface{}{"a": reflect.Value{}}},

		{Script: `func aFunc() {return a}; aFunc()`, Input: map[string]interface{}{"a": nil}, RunOutput: nil, Output: map[string]interface{}{"a": nil}},
		{Script: `func aFunc() {return a}; aFunc()`, Input: map[string]interface{}{"a": true}, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `func aFunc() {return a}; aFunc()`, Input: map[string]interface{}{"a": int64(1)}, RunOutput: int64(1), Output: map[string]interface{}{"a": int64(1)}},
		{Script: `func aFunc() {return a}; aFunc()`, Input: map[string]interface{}{"a": float64(1.1)}, RunOutput: float64(1.1), Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `func aFunc() {return a}; aFunc()`, Input: map[string]interface{}{"a": "a"}, RunOutput: "a", Output: map[string]interface{}{"a": "a"}},

		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": reflect.Value{}}, RunOutput: []interface{}{reflect.Value{}, reflect.Value{}}, Output: map[string]interface{}{"a": reflect.Value{}}},
		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": nil}, RunOutput: []interface{}{nil, nil}, Output: map[string]interface{}{"a": nil}},
		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": true}, RunOutput: []interface{}{true, true}, Output: map[string]interface{}{"a": true}},
		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": int32(1)}, RunOutput: []interface{}{int32(1), int32(1)}, Output: map[string]interface{}{"a": int32(1)}},
		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": int64(1)}, RunOutput: []interface{}{int64(1), int64(1)}, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": float32(1.1)}, RunOutput: []interface{}{float32(1.1), float32(1.1)}, Output: map[string]interface{}{"a": float32(1.1)}},
		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": float64(1.1)}, RunOutput: []interface{}{float64(1.1), float64(1.1)}, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `func aFunc() {return a, a}; aFunc()`, Input: map[string]interface{}{"a": "a"}, RunOutput: []interface{}{"a", "a"}, Output: map[string]interface{}{"a": "a"}},

		{Script: `func a(x) { return x}; a(nil)`, RunOutput: nil},
		{Script: `func a(x) { return x}; a(true)`, RunOutput: true},
		{Script: `func a(x) { return x}; a(1)`, RunOutput: int64(1)},
		{Script: `func a(x) { return x}; a(1.1)`, RunOutput: float64(1.1)},
		{Script: `func a(x) { return x}; a("a")`, RunOutput: "a"},

		{Script: `func aFunc() {return a}; for {aFunc(); break}`, Input: map[string]interface{}{"a": nil}, RunOutput: nil, Output: map[string]interface{}{"a": nil}},
		{Script: `func aFunc() {return a}; for {aFunc(); break}`, Input: map[string]interface{}{"a": true}, RunOutput: nil, Output: map[string]interface{}{"a": true}},
		{Script: `func aFunc() {return a}; for {aFunc(); break}`, Input: map[string]interface{}{"a": int64(1)}, RunOutput: nil, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `func aFunc() {return a}; for {aFunc(); break}`, Input: map[string]interface{}{"a": float64(1.1)}, RunOutput: nil, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `func aFunc() {return a}; for {aFunc(); break}`, Input: map[string]interface{}{"a": "a"}, RunOutput: nil, Output: map[string]interface{}{"a": "a"}},

		{Script: `func aFunc() {for {return a}}; aFunc()`, Input: map[string]interface{}{"a": nil}, RunOutput: nil, Output: map[string]interface{}{"a": nil}},
		{Script: `func aFunc() {for {return a}}; aFunc()`, Input: map[string]interface{}{"a": true}, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `func aFunc() {for {return a}}; aFunc()`, Input: map[string]interface{}{"a": int64(1)}, RunOutput: int64(1), Output: map[string]interface{}{"a": int64(1)}},
		{Script: `func aFunc() {for {return a}}; aFunc()`, Input: map[string]interface{}{"a": float64(1.1)}, RunOutput: float64(1.1), Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `func aFunc() {for {return a}}; aFunc()`, Input: map[string]interface{}{"a": "a"}, RunOutput: "a", Output: map[string]interface{}{"a": "a"}},

		{Script: `func aFunc() {for {if true {return a}}}; aFunc()`, Input: map[string]interface{}{"a": nil}, RunOutput: nil, Output: map[string]interface{}{"a": nil}},
		{Script: `func aFunc() {for {if true {return a}}}; aFunc()`, Input: map[string]interface{}{"a": true}, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `func aFunc() {for {if true {return a}}}; aFunc()`, Input: map[string]interface{}{"a": int64(1)}, RunOutput: int64(1), Output: map[string]interface{}{"a": int64(1)}},
		{Script: `func aFunc() {for {if true {return a}}}; aFunc()`, Input: map[string]interface{}{"a": float64(1.1)}, RunOutput: float64(1.1), Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `func aFunc() {for {if true {return a}}}; aFunc()`, Input: map[string]interface{}{"a": "a"}, RunOutput: "a", Output: map[string]interface{}{"a": "a"}},

		{Script: `func aFunc() {return nil, nil}; a, b = aFunc()`, RunOutput: nil, Output: map[string]interface{}{"a": nil, "b": nil}},
		{Script: `func aFunc() {return true, false}; a, b = aFunc()`, RunOutput: false, Output: map[string]interface{}{"a": true, "b": false}},
		{Script: `func aFunc() {return 1, 2}; a, b = aFunc()`, RunOutput: int64(2), Output: map[string]interface{}{"a": int64(1), "b": int64(2)}},
		{Script: `func aFunc() {return 1.1, 2.2}; a, b = aFunc()`, RunOutput: float64(2.2), Output: map[string]interface{}{"a": float64(1.1), "b": float64(2.2)}},
		{Script: `func aFunc() {return "a", "b"}; a, b = aFunc()`, RunOutput: "b", Output: map[string]interface{}{"a": "a", "b": "b"}},
	}
	testlib.Run(t, tests, nil)
}

func TestFunctions(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testlib.Test{
		{Script: `a()`, Input: map[string]interface{}{"a": reflect.Value{}}, RunError: fmt.Errorf("cannot call type reflect.Value")},
		{Script: `a = nil; a()`, RunError: fmt.Errorf("cannot call type interface {}"), Output: map[string]interface{}{"a": nil}},
		{Script: `a = true; a()`, RunError: fmt.Errorf("cannot call type bool"), Output: map[string]interface{}{"a": true}},
		{Script: `a = nil; b = func c(d) { return d == nil }; c = nil; c(a)`, RunError: fmt.Errorf("cannot call type interface {}"), Output: map[string]interface{}{"a": nil}},
		{Script: `a = [true]; a()`, RunError: fmt.Errorf("cannot call type []interface {}")},
		{Script: `a = [true]; func b(c) { return c() }; b(a)`, RunError: fmt.Errorf("cannot call type []interface {}")},
		{Script: `a = {}; a.missing()`, RunError: fmt.Errorf("cannot call type interface {}"), Output: map[string]interface{}{"a": map[interface{}]interface{}{}}},
		{Script: `a = 1; b = func(,a){}; a`, ParseError: fmt.Errorf("syntax error: unexpected ','"), RunOutput: int64(1)},

		{Script: `func a(b) { }; a()`, RunError: fmt.Errorf("function wants 1 arguments but received 0")},
		{Script: `func a(b) { }; a(true, true)`, RunError: fmt.Errorf("function wants 1 arguments but received 2")},
		{Script: `func a(b, c) { }; a()`, RunError: fmt.Errorf("function wants 2 arguments but received 0")},
		{Script: `func a(b, c) { }; a(true)`, RunError: fmt.Errorf("function wants 2 arguments but received 1")},
		{Script: `func a(b, c) { }; a(true, true, true)`, RunError: fmt.Errorf("function wants 2 arguments but received 3")},

		{Script: `func a() { return "a" }; a.b()`, RunError: fmt.Errorf("type func does not support member operation")},
		{Script: `a = [func () { return nil}]; func b(c) { return c() }; b(a[1])`, RunError: fmt.Errorf("index out of range")},
		{Script: `func a() { return "a" }; b()`, RunError: fmt.Errorf("undefined symbol 'b'")},
		{Script: ` func a() { return "a" }; 1++()`, RunError: fmt.Errorf("invalid operation")},
		{Script: ` func a(b) { return b }; a(1++)`, RunError: fmt.Errorf("invalid operation")},

		{Script: `a`, Input: map[string]interface{}{"a": testVarFunc}, RunOutput: testVarFunc, Output: map[string]interface{}{"a": testVarFunc}},
		{Script: `a()`, Input: map[string]interface{}{"a": testVarFunc}, RunOutput: int64(1), Output: map[string]interface{}{"a": testVarFunc}},
		{Script: `a`, Input: map[string]interface{}{"a": testVarFuncP}, RunOutput: testVarFuncP, Output: map[string]interface{}{"a": testVarFuncP}},
		// TOFIX:
		// {Script: `a()`, Input: map[string]interface{}{"a": testVarFuncP}, RunOutput: int64(1), Output: map[string]interface{}{"a": testVarFuncP}},

		{Script: `module a { func b() { return } }; a.b()`, RunOutput: nil},
		{Script: `module a { func b() { return nil} }; a.b()`, RunOutput: nil},
		{Script: `module a { func b() { return true} }; a.b()`, RunOutput: true},
		{Script: `module a { func b() { return 1} }; a.b()`, RunOutput: int64(1)},
		{Script: `module a { func b() { return 1.1} }; a.b()`, RunOutput: float64(1.1)},
		{Script: `module a { func b() { return "a"} }; a.b()`, RunOutput: "a"},

		{Script: `if true { module a { func b() { return } } }; a.b()`, RunOutput: nil},
		{Script: `if true { module a { func b() { return nil} } }; a.b()`, RunOutput: nil},
		{Script: `if true { module a { func b() { return true} } }; a.b()`, RunOutput: true},
		{Script: `if true { module a { func b() { return 1} } }; a.b()`, RunOutput: int64(1)},
		{Script: `if true { module a { func b() { return 1.1} } }; a.b()`, RunOutput: float64(1.1)},
		{Script: `if true { module a { func b() { return "a"} } }; a.b()`, RunOutput: "a"},

		{Script: `if true { module a { func b() { return 1} } }; a.b()`, RunOutput: int64(1)},

		{Script: `a = 1; func b() { a = 2 }; b()`, RunOutput: int64(2), Output: map[string]interface{}{"a": int64(2)}},
		{Script: `b(a); a`, Input: map[string]interface{}{"a": int64(1), "b": func(c interface{}) { c = int64(2); _ = c }}, RunOutput: int64(1), Output: map[string]interface{}{"a": int64(1)}},
		{Script: `func b() { }; go b()`, RunOutput: nil},

		{Script: `b(a)`, Input: map[string]interface{}{"a": nil, "b": func(c interface{}) bool { return c == nil }}, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": true, "b": func(c bool) bool { return c == true }}, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": int32(1), "b": func(c int32) bool { return c == 1 }}, RunOutput: true, Output: map[string]interface{}{"a": int32(1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": int64(1), "b": func(c int64) bool { return c == 1 }}, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": float32(1.1), "b": func(c float32) bool { return c == 1.1 }}, RunOutput: true, Output: map[string]interface{}{"a": float32(1.1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": float64(1.1), "b": func(c float64) bool { return c == 1.1 }}, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": "a", "b": func(c string) bool { return c == "a" }}, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueBool, "b": func(c reflect.Value) bool { return c == testVarValueBool }}, RunOutput: true, Output: map[string]interface{}{"a": testVarValueBool}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueInt32, "b": func(c reflect.Value) bool { return c == testVarValueInt32 }}, RunOutput: true, Output: map[string]interface{}{"a": testVarValueInt32}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueInt64, "b": func(c reflect.Value) bool { return c == testVarValueInt64 }}, RunOutput: true, Output: map[string]interface{}{"a": testVarValueInt64}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueFloat32, "b": func(c reflect.Value) bool { return c == testVarValueFloat32 }}, RunOutput: true, Output: map[string]interface{}{"a": testVarValueFloat32}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueFloat64, "b": func(c reflect.Value) bool { return c == testVarValueFloat64 }}, RunOutput: true, Output: map[string]interface{}{"a": testVarValueFloat64}},
		{Script: `b(a)`, Input: map[string]interface{}{"a": testVarValueString, "b": func(c reflect.Value) bool { return c == testVarValueString }}, RunOutput: true, Output: map[string]interface{}{"a": testVarValueString}},

		{Script: `x(a, b, c, d, e, f, g)`, Input: map[string]interface{}{"a": nil, "b": true, "c": int32(1), "d": int64(2), "e": float32(1.1), "f": float64(2.2), "g": "g",
			"x": func(a interface{}, b bool, c int32, d int64, e float32, f float64, g string) bool {
				return a == nil && b == true && c == 1 && d == 2 && e == 1.1 && f == 2.2 && g == "g"
			}}, RunOutput: true, Output: map[string]interface{}{"a": nil, "b": true, "c": int32(1), "d": int64(2), "e": float32(1.1), "f": float64(2.2), "g": "g"}},
		{Script: `x(a, b, c, d, e, f, g)`, Input: map[string]interface{}{"a": nil, "b": true, "c": int32(1), "d": int64(2), "e": float32(1.1), "f": float64(2.2), "g": "g",
			"x": func(a interface{}, b bool, c int32, d int64, e float32, f float64, g string) (interface{}, bool, int32, int64, float32, float64, string) {
				return a, b, c, d, e, f, g
			}}, RunOutput: []interface{}{nil, true, int32(1), int64(2), float32(1.1), float64(2.2), "g"}, Output: map[string]interface{}{"a": nil, "b": true, "c": int32(1), "d": int64(2), "e": float32(1.1), "f": float64(2.2), "g": "g"}},

		{Script: `a = nil; b(a)`, Input: map[string]interface{}{"b": func(c interface{}) bool { return c == nil }}, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `a = true; b(a)`, Input: map[string]interface{}{"b": func(c bool) bool { return c == true }}, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `a = 1; b(a)`, Input: map[string]interface{}{"b": func(c int64) bool { return c == 1 }}, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `a = 1.1; b(a)`, Input: map[string]interface{}{"b": func(c float64) bool { return c == 1.1 }}, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `a = "a"; b(a)`, Input: map[string]interface{}{"b": func(c string) bool { return c == "a" }}, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `func b(c) { return c == nil }; b(a)`, Input: map[string]interface{}{"a": nil}, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `func b(c) { return c == true }; b(a)`, Input: map[string]interface{}{"a": true}, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `func b(c) { return c == 1 }; b(a)`, Input: map[string]interface{}{"a": int32(1)}, RunOutput: true, Output: map[string]interface{}{"a": int32(1)}},
		{Script: `func b(c) { return c == 1 }; b(a)`, Input: map[string]interface{}{"a": int64(1)}, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `func b(c) { return c == 1.1 }; b(a)`, Input: map[string]interface{}{"a": float32(1.1)}, RunOutput: true, Output: map[string]interface{}{"a": float32(1.1)}},
		{Script: `func b(c) { return c == 1.1 }; b(a)`, Input: map[string]interface{}{"a": float64(1.1)}, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `func b(c) { return c == "a" }; b(a)`, Input: map[string]interface{}{"a": "a"}, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `a = nil; func b(c) { return c == nil }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `a = true; func b(c) { return c == true }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `a = 1; func b(c) { return c == 1 }; b(a)`, Input: map[string]interface{}{"a": int64(1)}, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `a = 1.1; func b(c) { return c == 1.1 }; b(a)`, Input: map[string]interface{}{"a": float64(1.1)}, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `a = "a"; func b(c) { return c == "a" }; b(a)`, Input: map[string]interface{}{"a": "a"}, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{nil}, "b": func(c interface{}) bool { return c == nil }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{nil}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{true}, "b": func(c interface{}) bool { return c == true }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{true}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{int32(1)}, "b": func(c interface{}) bool { return c == int32(1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{int32(1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{int64(1)}, "b": func(c interface{}) bool { return c == int64(1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{int64(1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{float32(1.1)}, "b": func(c interface{}) bool { return c == float32(1.1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{float32(1.1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{float64(1.1)}, "b": func(c interface{}) bool { return c == float64(1.1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{float64(1.1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{"a"}, "b": func(c interface{}) bool { return c == "a" }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{"a"}}},

		// TOFIX:
		//		{Script: `b(a)`,
		//			Input:     map[string]interface{}{"a": []bool{true, false, true}, "b": func(c ...bool) bool { return c[len(c)-1] }},
		//			RunOutput: true, Output: map[string]interface{}{"a": true}},

		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{true}, "b": func(c bool) bool { return c == true }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{true}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{int32(1)}, "b": func(c int32) bool { return c == int32(1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{int32(1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{int64(1)}, "b": func(c int64) bool { return c == int64(1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{int64(1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{float32(1.1)}, "b": func(c float32) bool { return c == float32(1.1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{float32(1.1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{float64(1.1)}, "b": func(c float64) bool { return c == float64(1.1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{float64(1.1)}}},
		{Script: `b(a[0])`, Input: map[string]interface{}{"a": []interface{}{"a"}, "b": func(c string) bool { return c == "a" }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{"a"}}},

		{Script: `a = [nil]; b(a[0])`, Input: map[string]interface{}{"b": func(c interface{}) bool { return c == nil }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{nil}}},
		{Script: `a = [true]; b(a[0])`, Input: map[string]interface{}{"b": func(c bool) bool { return c == true }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{true}}},
		{Script: `a = [1]; b(a[0])`, Input: map[string]interface{}{"b": func(c int64) bool { return c == int64(1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{int64(1)}}},
		{Script: `a = [1.1]; b(a[0])`, Input: map[string]interface{}{"b": func(c float64) bool { return c == float64(1.1) }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{float64(1.1)}}},
		{Script: `a = ["a"]; b(a[0])`, Input: map[string]interface{}{"b": func(c string) bool { return c == "a" }}, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{"a"}}},

		{Script: `a = [nil]; func b(c) { c == nil }; b(a[0])`, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{nil}}},
		{Script: `a = [true]; func b(c) { c == true }; b(a[0])`, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{true}}},
		{Script: `a = [1]; func b(c) { c == 1 }; b(a[0])`, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{int64(1)}}},
		{Script: `a = [1.1]; func b(c) { c == 1.1 }; b(a[0])`, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{float64(1.1)}}},
		{Script: `a = ["a"]; func b(c) { c == "a" }; b(a[0])`, RunOutput: true, Output: map[string]interface{}{"a": []interface{}{"a"}}},

		{Script: `a = nil; b = func (d) { return d == nil }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `a = true; b = func (d) { return d == true }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `a = 1; b = func (d) { return d == 1 }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `a = 1.1; b = func (d) { return d == 1.1 }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `a = "a"; b = func (d) { return d == "a" }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `a = nil; b = func c(d) { return d == nil }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `a = true; b = func c(d) { return d == true }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `a = 1; b = func c(d) { return d == 1 }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `a = 1.1; b = func c(d) { return d == 1.1 }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `a = "a"; b = func c(d) { return d == "a" }; b(a)`, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `a = nil; b = func c(d) { return d == nil }; c(a)`, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `a = true; b = func c(d) { return d == true }; c(a)`, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `a = 1; b = func c(d) { return d == 1 }; c(a)`, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `a = 1.1; b = func c(d) { return d == 1.1 }; c(a)`, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `a = "a"; b = func c(d) { return d == "a" }; c(a)`, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `a = nil; func b() { return func c(d) { d == nil } }; e = b(); e(a)`, RunOutput: true, Output: map[string]interface{}{"a": nil}},
		{Script: `a = true; func b() { return func c(d) { d == true } }; e = b(); e(a)`, RunOutput: true, Output: map[string]interface{}{"a": true}},
		{Script: `a = 1; func b() { return func c(d) { d == 1 } }; e = b(); e(a)`, RunOutput: true, Output: map[string]interface{}{"a": int64(1)}},
		{Script: `a = 1.1; func b() { return func c(d) { d == 1.1 } }; e = b(); e(a)`, RunOutput: true, Output: map[string]interface{}{"a": float64(1.1)}},
		{Script: `a = "a"; func b() { return func c(d) { d == "a" } }; e = b(); e(a)`, RunOutput: true, Output: map[string]interface{}{"a": "a"}},

		{Script: `a = func () { return nil }; func b(c) { return c() }; b(a)`, RunOutput: nil},
		{Script: `a = func () { return true }; func b(c) { return c() }; b(a)`, RunOutput: true},
		{Script: `a = func () { return 1 }; func b(c) { return c() }; b(a)`, RunOutput: int64(1)},
		{Script: `a = func () { return 1.1 }; func b(c) { return c() }; b(a)`, RunOutput: float64(1.1)},
		{Script: `a = func () { return "a" }; func b(c) { return c() }; b(a)`, RunOutput: "a"},

		{Script: `a = [nil]; func c(d) { return d[0] }; c(a)`, RunOutput: nil},
		{Script: `a = [true]; func c(d) { return d[0] }; c(a)`, RunOutput: true},
		{Script: `a = [1]; func c(d) { return d[0] }; c(a)`, RunOutput: int64(1)},
		{Script: `a = [1.1]; func c(d) { return d[0] }; c(a)`, RunOutput: float64(1.1)},
		{Script: `a = ["a"]; func c(d) { return d[0] }; c(a)`, RunOutput: "a"},

		{Script: `a = {"b": nil}; func c(d) { return d.b }; c(a)`, RunOutput: nil},
		{Script: `a = {"b": true}; func c(d) { return d.b }; c(a)`, RunOutput: true},
		{Script: `a = {"b": 1}; func c(d) { return d.b }; c(a)`, RunOutput: int64(1)},
		{Script: `a = {"b": 1.1}; func c(d) { return d.b }; c(a)`, RunOutput: float64(1.1)},
		{Script: `a = {"b": "a"}; func c(d) { return d.b }; c(a)`, RunOutput: "a"},

		{Script: `a = func() { return func(c) { return c + "c"} }; a()("a")`, RunOutput: "ac"},
		{Script: `a = func() { return func(c) { return c + "c"} }(); a("a")`, RunOutput: "ac"},
		{Script: `a = func() { return func(c) { return c + "c"} }()("a")`, RunOutput: "ac"},
		{Script: `func() { return func(c) { return c + "c"} }()("a")`, RunOutput: "ac"},

		{Script: `a = func(b) { return func() { return b + "c"} }; b = a("a"); b()`, RunOutput: "ac"},
		{Script: `a = func(b) { return func() { return b + "c"} }("a"); a()`, RunOutput: "ac"},
		{Script: `a = func(b) { return func() { return b + "c"} }("a")()`, RunOutput: "ac"},
		{Script: `func(b) { return func() { return b + "c"} }("a")()`, RunOutput: "ac"},

		{Script: `a = func(b) { return func(c) { return b[c] } }; b = a({"x": "x"}); b("x")`, RunOutput: "x"},
		{Script: `a = func(b) { return func(c) { return b[c] } }({"x": "x"}); a("x")`, RunOutput: "x"},
		{Script: `a = func(b) { return func(c) { return b[c] } }({"x": "x"})("x")`, RunOutput: "x"},
		{Script: `func(b) { return func(c) { return b[c] } }({"x": "x"})("x")`, RunOutput: "x"},

		{Script: `a = func(b) { return func(c) { return b[c] } }; x = {"y": "y"}; b = a(x); x = {"y": "y"}; b("y")`, RunOutput: "y"},
		{Script: `a = func(b) { return func(c) { return b[c] } }; x = {"y": "y"}; b = a(x); x.y = "z"; b("y")`, RunOutput: "z"},

		{Script: ` func a() { return "a" }; a()`, RunOutput: "a"},
		{Script: `a = func a() { return "a" }; a = func() { return "b" }; a()`, RunOutput: "b"},
		{Script: `a = "a.b"; func a() { return "a" }; a()`, RunOutput: "a"},

		{Script: `a = func() { b = "b"; return func() { b += "c" } }(); a()`, RunOutput: "bc"},
		{Script: `a = func() { b = "b"; return func() { b += "c"; return b} }(); a()`, RunOutput: "bc"},
		{Script: `a = func(b) { return func(c) { return func(d) { return d + "d" }(c) + "c" }(b) + "b" }("a")`, RunOutput: "adcb"},
		{Script: `a = func(b) { return "b" + func(c) { return "c" + func(d) { return "d" + d  }(c) }(b) }("a")`, RunOutput: "bcda"},
		{Script: `a = func(b) { return b + "b" }; a( func(c) { return c + "c" }("a") )`, RunOutput: "acb"},

		{Script: `a = func(x, y) { return func() { x(y) } }; b = a(func (z) { return z + "z" }, "b"); b()`, RunOutput: "bz"},

		{Script: `a = make(Time); a.IsZero()`, Types: map[string]interface{}{"Time": time.Time{}}, RunOutput: true},

		{Script: `a = make(Buffer); n, err = a.WriteString("a"); if err != nil { return err }; n`, Types: map[string]interface{}{"Buffer": bytes.Buffer{}}, RunOutput: 1},
		{Script: `a = make(Buffer); n, err = a.WriteString("a"); if err != nil { return err }; a.String()`, Types: map[string]interface{}{"Buffer": bytes.Buffer{}}, RunOutput: "a"},
	}
	testlib.Run(t, tests, nil)
}

func TestPointerFunctions(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	testFunctionPointer := func(b interface{}) string {
		rv := reflect.ValueOf(b)
		if !rv.IsValid() {
			return "invalid"
		}
		if rv.Kind() != reflect.Ptr {
			return fmt.Sprintf("not ptr: " + rv.Kind().String())
		}
		if rv.IsNil() {
			return "IsNil"
		}
		if !rv.Elem().CanInterface() {
			return "cannot interface"
		}
		if rv.Elem().Interface() != int64(1) {
			return fmt.Sprintf("not 1: %v", rv.Elem().Interface())
		}
		if !rv.Elem().CanSet() {
			return "cannot set"
		}
		slice := reflect.MakeSlice(interfaceSliceType, 0, 1)
		value, _ := makeValue(stringType)
		value.SetString("b")
		slice = reflect.Append(slice, value)
		rv.Elem().Set(slice)
		return "good"
	}
	tests := []testlib.Test{
		{Script: `b = 1; a(&b)`, Input: map[string]interface{}{"a": testFunctionPointer}, RunOutput: "good", Output: map[string]interface{}{"b": []interface{}{"b"}}},
	}
	testlib.Run(t, tests, nil)
}

func TestVariadicFunctions(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testlib.Test{
		// params Variadic arg !Variadic
		{Script: `func a(b...) { return b }; a()`, RunOutput: []interface{}{}},
		{Script: `func a(b...) { return b }; a(true)`, RunOutput: []interface{}{true}},
		{Script: `func a(b...) { return b }; a(true, true)`, RunOutput: []interface{}{true, true}},
		{Script: `func a(b...) { return b }; a([true])`, RunOutput: []interface{}{[]interface{}{true}}},
		{Script: `func a(b...) { return b }; a([true, true])`, RunOutput: []interface{}{[]interface{}{true, true}}},
		{Script: `func a(b...) { return b }; a([true, true], [true, true])`, RunOutput: []interface{}{[]interface{}{true, true}, []interface{}{true, true}}},

		// params Variadic arg !Variadic
		{Script: `func a(b, c...) { return c }; a()`, RunError: fmt.Errorf("function wants 2 arguments but received 0")},
		{Script: `func a(b, c...) { return c }; a(true)`, RunOutput: []interface{}{}},
		{Script: `func a(b, c...) { return c }; a(true, true)`, RunOutput: []interface{}{true}},
		{Script: `func a(b, c...) { return c }; a(true, true, true)`, RunOutput: []interface{}{true, true}},
		{Script: `func a(b, c...) { return c }; a([true])`, RunOutput: []interface{}{}},
		{Script: `func a(b, c...) { return c }; a([true], [true])`, RunOutput: []interface{}{[]interface{}{true}}},
		{Script: `func a(b, c...) { return c }; a([true], [true], [true])`, RunOutput: []interface{}{[]interface{}{true}, []interface{}{true}}},
		{Script: `func a(b, c...) { return c }; a([true], [true, true], [true, true])`, RunOutput: []interface{}{[]interface{}{true, true}, []interface{}{true, true}}},

		// params Variadic arg Variadic
		{Script: `func a(b...) { return b }; a([true]...)`, RunOutput: []interface{}{true}},
		{Script: `func a(b...) { return b }; a([true, true]...)`, RunOutput: []interface{}{true, true}},
		{Script: `func a(b...) { return b }; a(true, [true]...)`, RunError: fmt.Errorf("function wants 1 arguments but received 2")},

		// params Variadic arg Variadic
		{Script: `func a(b, c...) { return c }; a([true]...)`, RunOutput: []interface{}{}},
		{Script: `func a(b, c...) { return c }; a([true, true]...)`, RunOutput: []interface{}{}},
		{Script: `func a(b, c...) { return c }; a(true, [true]...)`, RunOutput: []interface{}{true}},
		{Script: `func a(b, c...) { return c }; a(true, [true, true]...)`, RunOutput: []interface{}{true, true}},

		// params !Variadic arg Variadic
		{Script: `func a() { return "a" }; a([true]...)`, RunOutput: "a"},
		{Script: `func a() { return "a" }; a(true, [true]...)`, RunOutput: "a"},
		{Script: `func a() { return "a" }; a(true, [true, true]...)`, RunOutput: "a"},

		// params !Variadic arg Variadic
		{Script: `func a(b) { return b }; a(true...)`, RunError: fmt.Errorf("call is variadic but last parameter is of type bool")},
		{Script: `func a(b) { return b }; a([true]...)`, RunOutput: true},
		{Script: `func a(b) { return b }; a(true, false...)`, RunError: fmt.Errorf("function wants 1 arguments but received 2")},
		{Script: `func a(b) { return b }; a(true, [1]...)`, RunError: fmt.Errorf("function wants 1 arguments but received 2")},
		{Script: `func a(b) { return b }; a(true, [1, 2]...)`, RunError: fmt.Errorf("function wants 1 arguments but received 2")},
		{Script: `func a(b) { return b }; a([true, 1]...)`, RunOutput: true},
		{Script: `func a(b) { return b }; a([true, 1, 2]...)`, RunOutput: true},

		// params !Variadic arg Variadi
		{Script: `func a(b, c) { return c }; a(false...)`, RunError: fmt.Errorf("call is variadic but last parameter is of type bool")},
		{Script: `func a(b, c) { return c }; a([1]...)`, RunError: fmt.Errorf("function wants 2 arguments but received 1")},
		{Script: `func a(b, c) { return c }; a(1, true...)`, RunError: fmt.Errorf("call is variadic but last parameter is of type bool")},
		{Script: `func a(b, c) { return c }; a(1, [true]...)`, RunOutput: true},
		{Script: `func a(b, c) { return c }; a([1, true]...)`, RunOutput: true},
		{Script: `func a(b, c) { return c }; a(1, true...)`, RunError: fmt.Errorf("call is variadic but last parameter is of type bool")},
		{Script: `func a(b, c) { return c }; a(1, [true]...)`, RunOutput: true},
		{Script: `func a(b, c) { return c }; a(1, true, false...)`, RunError: fmt.Errorf("function wants 2 arguments but received 3")},
		{Script: `func a(b, c) { return c }; a(1, true, [2]...)`, RunError: fmt.Errorf("function wants 2 arguments but received 3")},
		{Script: `func a(b, c) { return c }; a(1, [true, 2]...)`, RunOutput: true},
		{Script: `func a(b, c) { return c }; a([1, true, 2]...)`, RunOutput: true},
	}
	testlib.Run(t, tests, nil)
}

func TestFunctionsInArraysAndMaps(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testlib.Test{
		{Script: `a = [func () { return nil }]; a[0]()`, RunOutput: nil},
		{Script: `a = [func () { return true }]; a[0]()`, RunOutput: true},
		{Script: `a = [func () { return 1 }]; a[0]()`, RunOutput: int64(1)},
		{Script: `a = [func () { return 1.1 }]; a[0]()`, RunOutput: float64(1.1)},
		{Script: `a = [func () { return "a" }]; a[0]()`, RunOutput: "a"},

		{Script: `a = [func () { return nil }]; b = a[0]; b()`, RunOutput: nil},
		{Script: `a = [func () { return true }]; b = a[0]; b()`, RunOutput: true},
		{Script: `a = [func () { return 1 }]; b = a[0]; b()`, RunOutput: int64(1)},
		{Script: `a = [func () { return 1.1 }]; b = a[0]; b()`, RunOutput: float64(1.1)},
		{Script: `a = [func () { return "a" }]; b = a[0]; b()`, RunOutput: "a"},

		{Script: `a = [func () { return nil}]; func b(c) { return c() }; b(a[0])`, RunOutput: nil},
		{Script: `a = [func () { return true}]; func b(c) { return c() }; b(a[0])`, RunOutput: true},
		{Script: `a = [func () { return 1}]; func b(c) { return c() }; b(a[0])`, RunOutput: int64(1)},
		{Script: `a = [func () { return 1.1}]; func b(c) { return c() }; b(a[0])`, RunOutput: float64(1.1)},
		{Script: `a = [func () { return "a"}]; func b(c) { return c() }; b(a[0])`, RunOutput: "a"},

		{Script: `a = {"b": func () { return nil }}; a["b"]()`, RunOutput: nil},
		{Script: `a = {"b": func () { return true }}; a["b"]()`, RunOutput: true},
		{Script: `a = {"b": func () { return 1 }}; a["b"]()`, RunOutput: int64(1)},
		{Script: `a = {"b": func () { return 1.1 }}; a["b"]()`, RunOutput: float64(1.1)},
		{Script: `a = {"b": func () { return "a" }}; a["b"]()`, RunOutput: "a"},

		{Script: `a = {"b": func () { return nil }}; a.b()`, RunOutput: nil},
		{Script: `a = {"b": func () { return true }}; a.b()`, RunOutput: true},
		{Script: `a = {"b": func () { return 1 }}; a.b()`, RunOutput: int64(1)},
		{Script: `a = {"b": func () { return 1.1 }}; a.b()`, RunOutput: float64(1.1)},
		{Script: `a = {"b": func () { return "a" }}; a.b()`, RunOutput: "a"},

		{Script: `a = {"b": func () { return nil }}; func c(d) { return d() }; c(a.b)`, RunOutput: nil},
		{Script: `a = {"b": func () { return true }}; func c(d) { return d() }; c(a.b)`, RunOutput: true},
		{Script: `a = {"b": func () { return 1 }}; func c(d) { return d() }; c(a.b)`, RunOutput: int64(1)},
		{Script: `a = {"b": func () { return 1.1 }}; func c(d) { return d() }; c(a.b)`, RunOutput: float64(1.1)},
		{Script: `a = {"b": func () { return "a" }}; func c(d) { return d() }; c(a.b)`, RunOutput: "a"},
	}
	testlib.Run(t, tests, nil)
}

func TestFunctionConvertions(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "0")
	tests := []testlib.Test{
		{Script: `b = func(c){ return c }; a("x", b)`, Input: map[string]interface{}{"a": func(b string, c func(string) string) string { return c(b) }}, RunOutput: "x"},
		{Script: `b = make(struct); b.A = func (c, d) { return c == d }; b.A(2, 2)`, Types: map[string]interface{}{"struct": &struct {
			A func(int, int) bool
		}{}},
			RunOutput: true},
		{Script: `b = 1; a(&b)`, Input: map[string]interface{}{"a": func(b *int64) { *b = int64(2) }}, Output: map[string]interface{}{"b": int64(2)}},
		// TOFIX:
		// {Script: `b = 1; c = &b; a(c); *c`, Input: map[string]interface{}{"a": func(b *int64) { *b = int64(2) }}, RunOutput: int64(2), Output: map[string]interface{}{"b": int64(2)}},
		// TODO: add more tests
	}
	testlib.Run(t, tests, nil)
}

func TestVariadicFunctionConvertions(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	testSumFunc := func(nums ...int64) int64 {
		var total int64
		for _, num := range nums {
			total += num
		}
		return total
	}
	tests := []testlib.Test{
		// params Variadic arg !Variadic
		{Script: `a(true)`, Input: map[string]interface{}{"a": func(b ...interface{}) []interface{} { return b }}, RunOutput: []interface{}{true}},

		{Script: `a()`, Input: map[string]interface{}{"a": testSumFunc}, RunOutput: int64(0)},
		{Script: `a(1)`, Input: map[string]interface{}{"a": testSumFunc}, RunOutput: int64(1)},
		{Script: `a(1, 2)`, Input: map[string]interface{}{"a": testSumFunc}, RunOutput: int64(3)},
		{Script: `a(1, 2, 3)`, Input: map[string]interface{}{"a": testSumFunc}, RunOutput: int64(6)},

		// TODO: add more tests
	}
	testlib.Run(t, tests, nil)
}

func TestLen(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testlib.Test{
		{Script: `len(1++)`, RunError: fmt.Errorf("invalid operation")},
		{Script: `len(true)`, RunError: fmt.Errorf("type bool does not support len operation")},

		{Script: `a = ""; len(a)`, RunOutput: int64(0)},
		{Script: `a = "test"; len(a)`, RunOutput: int64(4)},
		{Script: `a = []; len(a)`, RunOutput: int64(0)},
		{Script: `a = [nil]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [true]; len(a)`, RunOutput: int64(1)},
		{Script: `a = ["test"]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [1]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [1.1]; len(a)`, RunOutput: int64(1)},

		{Script: `a = [[]]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [[nil]]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [[true]]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [["test"]]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [[1]]; len(a)`, RunOutput: int64(1)},
		{Script: `a = [[1.1]]; len(a)`, RunOutput: int64(1)},

		{Script: `a = [[]]; len(a[0])`, RunOutput: int64(0)},
		{Script: `a = [[nil]]; len(a[0])`, RunOutput: int64(1)},
		{Script: `a = [[true]]; len(a[0])`, RunOutput: int64(1)},
		{Script: `a = [["test"]]; len(a[0])`, RunOutput: int64(1)},
		{Script: `a = [[1]]; len(a[0])`, RunOutput: int64(1)},
		{Script: `a = [[1.1]]; len(a[0])`, RunOutput: int64(1)},

		{Script: `len(a)`, Input: map[string]interface{}{"a": "a"}, RunOutput: int64(1), Output: map[string]interface{}{"a": "a"}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": map[string]interface{}{}}, RunOutput: int64(0), Output: map[string]interface{}{"a": map[string]interface{}{}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": map[string]interface{}{"test": "test"}}, RunOutput: int64(1), Output: map[string]interface{}{"a": map[string]interface{}{"test": "test"}}},
		{Script: `len(a["test"])`, Input: map[string]interface{}{"a": map[string]interface{}{"test": "test"}}, RunOutput: int64(4), Output: map[string]interface{}{"a": map[string]interface{}{"test": "test"}}},

		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{}}, RunOutput: int64(0), Output: map[string]interface{}{"a": []interface{}{}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{nil}}, RunOutput: int64(1), Output: map[string]interface{}{"a": []interface{}{nil}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{true}}, RunOutput: int64(1), Output: map[string]interface{}{"a": []interface{}{true}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{int32(1)}}, RunOutput: int64(1), Output: map[string]interface{}{"a": []interface{}{int32(1)}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{int64(1)}}, RunOutput: int64(1), Output: map[string]interface{}{"a": []interface{}{int64(1)}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{float32(1.1)}}, RunOutput: int64(1), Output: map[string]interface{}{"a": []interface{}{float32(1.1)}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{float64(1.1)}}, RunOutput: int64(1), Output: map[string]interface{}{"a": []interface{}{float64(1.1)}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": []interface{}{"a"}}, RunOutput: int64(1), Output: map[string]interface{}{"a": []interface{}{"a"}}},

		{Script: `len(a[0])`, Input: map[string]interface{}{"a": []interface{}{"test"}}, RunOutput: int64(4), Output: map[string]interface{}{"a": []interface{}{"test"}}},

		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{}}, RunOutput: int64(0), Output: map[string]interface{}{"a": [][]interface{}{}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{nil}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{nil}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{{nil}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{nil}}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{{true}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{true}}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{{int32(1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{int32(1)}}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{{int64(1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{int64(1)}}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{{float32(1.1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{float32(1.1)}}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{{float64(1.1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{float64(1.1)}}}},
		{Script: `len(a)`, Input: map[string]interface{}{"a": [][]interface{}{{"a"}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{"a"}}}},

		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{nil}}, RunOutput: int64(0), Output: map[string]interface{}{"a": [][]interface{}{nil}}},
		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{{nil}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{nil}}}},
		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{{true}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{true}}}},
		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{{int32(1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{int32(1)}}}},
		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{{int64(1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{int64(1)}}}},
		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{{float32(1.1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{float32(1.1)}}}},
		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{{float64(1.1)}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{float64(1.1)}}}},
		{Script: `len(a[0])`, Input: map[string]interface{}{"a": [][]interface{}{{"a"}}}, RunOutput: int64(1), Output: map[string]interface{}{"a": [][]interface{}{{"a"}}}},

		{Script: `len(a[0][0])`, Input: map[string]interface{}{"a": [][]interface{}{{"test"}}}, RunOutput: int64(4), Output: map[string]interface{}{"a": [][]interface{}{{"test"}}}},
	}
	testlib.Run(t, tests, nil)
}

func TestCallFunctionWithVararg(t *testing.T) {
	env := NewEnv()
	err := env.Define("X", func(args ...string) []string {
		return args
	})
	if err != nil {
		t.Errorf("Define error: %v", err)
	}
	want := []string{"foo", "bar", "baz"}
	err = env.Define("a", want)
	if err != nil {
		t.Errorf("Define error: %v", err)
	}
	got, err := env.Execute(`X(a...)`)
	if err != nil {
		t.Errorf("execute error - received %#v - expected: %#v", err, ErrInterrupt)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("execute error - received %#v - expected: %#v", got, want)
	}
}
