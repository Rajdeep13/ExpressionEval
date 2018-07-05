package datatype

import (
	"fmt"
	"time"
)

// Type conversion helper functions

func ToBool(d DataType) (Bool, error) {
	if v, ok := d.(interface {
		ToBool() (Bool, error)
	}); ok {
		return v.ToBool()
	} else {
		return false, fmt.Errorf("Cannot convert '%v' to Bool", d)
	}
}
func ToDateTime(d DataType) (DateTime, error) {
	if v, ok := d.(interface {
		ToDateTime() (DateTime, error)
	}); ok {
		return v.ToDateTime()
	} else {
		return DateTime(time.Now()), fmt.Errorf("Cannot convert '%v' to DateTime", d)
	}
}
func ToInt(d DataType) (Int, error) {
	if v, ok := d.(interface {
		ToInt() (Int, error)
	}); ok {
		return v.ToInt()
	} else {
		return 0, fmt.Errorf("Cannot convert '%v' to Int", d)
	}
}
func ToDouble(d DataType) (Double, error) {
	if v, ok := d.(interface {
		ToDouble() (Double, error)
	}); ok {
		return v.ToDouble()
	} else {
		return 0, fmt.Errorf("Cannot convert '%v' to Double", d)
	}
}
func ToString(d DataType) (String, error) {
	if v, ok := d.(interface {
		ToString() (String, error)
	}); ok {
		return v.ToString()
	} else {
		return "", fmt.Errorf("Cannot convert '%v' to String", d)
	}
}
func ToChar(d DataType) (Char, error) {
	if v, ok := d.(interface {
		ToChar() (Char, error)
	}); ok {
		return v.ToChar()
	} else {
		return '?', fmt.Errorf("Cannot convert '%v' to Char", d)
	}
}
