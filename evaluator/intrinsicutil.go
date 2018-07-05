package evaluator

import (
	"fmt"
	"strings"

	"github.com/contactkeval/expressioneval/datatype"
	"github.com/davecgh/go-spew/spew"
)

// An intrinsic method
type intrinsicMethod interface {
	ExecuteMethod(pfe *PostfixExpression) (datatype.DataType, error)
}

// Simple wrapper over a func to convert it to an intrinsicMethod
type intrinsicMethodFunc func(pfe *PostfixExpression) (datatype.DataType, error)

func (imf intrinsicMethodFunc) ExecuteMethod(pfe *PostfixExpression) (datatype.DataType, error) {
	return imf(pfe)
}

// List of alternative intrinsic methods. The result is anoter intrinsicMethod
// that calls each of the methods in the list till one of them succeeds.
type intrinsicMethodList []intrinsicMethod

func (iml intrinsicMethodList) ExecuteMethod(pfe *PostfixExpression) (datatype.DataType, error) {
	var errs error
	for _, im := range iml {
		if d, err := im.ExecuteMethod(pfe); err != nil {
			if errs != nil {
				errs = fmt.Errorf("%s\n%s", errs.Error(), err.Error())
			} else {
				errs = err
			}
		} else {
			return d, nil
		}
	}
	return nil, errs
}

// Type conversion utilities

func toInt(d datatype.DataType) int {
	num, _ := datatype.ToInt(d)
	return int(num)
}

func toFloat(d datatype.DataType) float64 {
	num, _ := datatype.ToDouble(d)
	return float64(num)
}

func toBool(d datatype.DataType) bool {
	v, _ := datatype.ToBool(d)
	return bool(v)
}

func toString(d datatype.DataType) string {
	v, _ := datatype.ToString(d)
	return string(v)
}

func toRune(d datatype.DataType) rune {
	v, _ := datatype.ToChar(d)
	return rune(v)
}

func toSlice(d datatype.DataType) []datatype.DataType {
	l, _ := d.(datatype.List)
	return []datatype.DataType(l)
}

func toBoolSlice(l datatype.List) []bool {
	ls := []bool{}
	for _, item := range l {
		v := toBool(item)
		ls = append(ls, v)
	}
	return ls
}
func toIntSlice(l datatype.List) []int {
	ls := []int{}
	for _, item := range l {
		v := toInt(item)
		ls = append(ls, v)
	}
	return ls
}
func toFloatSlice(l datatype.List) []float64 {
	ls := []float64{}
	for _, item := range l {
		v := toFloat(item)
		ls = append(ls, v)
	}
	return ls
}
func toRuneSlice(l datatype.List) []rune {
	ls := []rune{}
	for _, item := range l {
		v := toRune(item)
		ls = append(ls, v)
	}
	return ls
}
func toStringSlice(l datatype.List) []string {
	ls := []string{}
	for _, item := range l {
		v := toString(item)
		ls = append(ls, v)
	}
	return ls
}

func GetIntrinsicMethod(name string) (intrinsicMethod, error) {
	if meth, exists := intrinsicMethods[name]; exists {
		return meth, nil
	} else {
		return nil, fmt.Errorf("Method %s does not exist", name)
	}
}

// Polymorphic methods, with some basic type checking
// Syntax: "types1", func1, "types2", func2, ...
//
// Supported types:
// S  - String
// N  - Number (Int or Double)
// I  - Int (I0, I1 are alternatives with a default value)
// D  - Double
// C  - Char
// B  - Bool (BF, BT are alternatives with a default value)
// L  - List
// LS - List of strings
// LN - List of Numbers
//
// Eg. S,LS,BF - func takes 3 arguments - string, list of strings and an
// optional bool with 'false' as the default value.
func polyTypeCheckedMethod(typesAndFuncs ...interface{}) intrinsicMethod {
	return intrinsicMethodFunc(func(pfe *PostfixExpression) (datatype.DataType, error) {
		// It's either a single operand or a single CommaList which contains a
		// list of operands
		funcArg, err := GetUnaryOperand(pfe)
		if err != nil {
			return nil, err
		}

		var args []datatype.DataType

		switch v := funcArg.(type) {
		case datatype.CommaList:
			args = []datatype.DataType(v.Flatten())
		default:
			args = []datatype.DataType{v}
		}

		unmatchedTypes := ""

		tf := typesAndFuncs
		spew.Dump(tf)
	Outer:
		for len(tf) >= 2 {
			typeString, f := tf[0].(string), tf[1].(func(args ...datatype.DataType) (datatype.DataType, error))
			tf = tf[2:]

			// only used in an error message if we break out of the loop
			unmatchedTypes = unmatchedTypes + " " + typeString

			types := strings.Split(typeString, ",")

			for i, arg := range args {
				// All possible types in the type signature
				switch types[i] {
				case "S":
					if _, ok := arg.(datatype.String); !ok {
						continue Outer
					}
				case "N":
					if !datatype.IsNumber(arg) {
						continue Outer
					}
				case "I", "I0", "I1":
					if _, ok := arg.(datatype.Int); !ok {
						continue Outer
					}
				case "D":
					if _, ok := arg.(datatype.Double); !ok {
						continue Outer
					}
				case "C":
					if _, ok := arg.(datatype.Char); !ok {
						continue Outer
					}
				case "B", "BT", "BF":
					if _, ok := arg.(datatype.Bool); !ok {
						continue Outer
					}
				case "L":
					if _, ok := arg.(datatype.List); !ok {
						continue Outer
					}
				case "LS": // List of strings
					if l, ok := arg.(datatype.List); !ok || !datatype.IsString(l...) {
						continue Outer
					}
				case "LN": // List of numbers
					if l, ok := arg.(datatype.List); !ok || !datatype.IsNumber(l...) {
						continue Outer
					}
				default:
					panic("Unknown type signature")
				}
			}

			// Special support for some default args
			if len(types) > len(args) {
				for i := len(args); i < len(types); i++ {
					switch types[i] {
					case "I0":
						args = append(args, datatype.Int(0))
					case "I1":
						args = append(args, datatype.Int(1))
					case "BT":
						args = append(args, datatype.Bool(true))
					case "BF":
						args = append(args, datatype.Bool(false))
					default:
						break
					}
				}
			}

			// This type signature did not match the number of arguments passed
			// in
			if len(args) != len(types) {
				continue
			}

			return f(args...)
		}

		return nil, fmt.Errorf("Intrinsic method only accepts arguments of signatures: %s", unmatchedTypes)
	})
}
