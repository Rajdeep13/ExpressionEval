package datatype

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Data type names. DataType tokens return one of these when their DataType()
// method is called.
const (
	DataTypeBool              = "Bool"
	DataTypeString            = "String"
	DataTypeDouble            = "Double"
	DataTypeInt               = "Int32"
	DataTypeChar              = "Char"
	DataTypeDateTime          = "DateTime"
	DataTypeIntRange          = "IntRange"
	DataTypeConditionalValues = "ConditionalValues"
)

// All values should be a DatatType
type DataType interface {
	DataType() string
}

// Possible data types, defined in terms of Go data types
type Char rune
type String string
type Int int
type Double float64
type Bool bool
type DateTime time.Time
type IntRange struct{ From, To int }
type ConditionalValues struct {
	Cond        Bool
	True, False DataType
}

// Null doesn't seem to have been implemented properly in the C# code
// type Null struct{}

// Data types can implement this interface to specify how they should be
// printed on screen
type Printable interface {
	ToPrint() string
}

/////////////////////////////
// Bool
func (b Bool) DataType() string { return DataTypeBool }
func (b Bool) ToBool() (Bool, error) {
	return b, nil
}
func (b Bool) ToString() (String, error) {
	if b {
		return "True", nil
	} else {
		return "False", nil
	}
}
func (b Bool) ToInt() (Int, error) {
	if b {
		return 1, nil
	} else {
		return 0, nil
	}
}
func (b Bool) ToDouble() (Double, error) {
	if b {
		return 1, nil
	} else {
		return 0, nil
	}
}

/////////////////////////////
// DateTime
func (b DateTime) DataType() string { return DataTypeDateTime }
func (n DateTime) ToDateTime() (DateTime, error) {
	return n, nil
}
func (n DateTime) ToString() (String, error) {
	return String(time.Time(n).Format("01/02/2006")), nil
}
func (n DateTime) ToPrint() string {
	return time.Time(n).Format("01/02/2006")
}

/////////////////////////////
// Int
func (b Int) DataType() string { return DataTypeInt }
func (n Int) ToInt() (Int, error) {
	return n, nil
}
func (n Int) ToBool() (Bool, error) {
	if n != 0 {
		return true, nil
	} else {
		return false, nil
	}
}
func (n Int) ToString() (String, error) {
	return String(strconv.Itoa(int(n))), nil
}
func (n Int) ToDouble() (Double, error) {
	return Double(n), nil
}

/////////////////////////////
// Double
func (b Double) DataType() string { return DataTypeDouble }
func (n Double) ToDouble() (Double, error) {
	return n, nil
}
func (n Double) ToBool() (Bool, error) {
	if n != 0 {
		return true, nil
	} else {
		return false, nil
	}
}
func (n Double) ToString() (String, error) {
	return String(fmt.Sprintf("%f", n)), nil
}
func (n Double) ToInt() (Int, error) {
	return Int(n), nil
}

/////////////////////////////
// Char
func (c Char) DataType() string { return DataTypeChar }
func (c Char) ToChar() (Char, error) {
	return c, nil
}
func (c Char) ToString() (String, error) {
	return String(rune(c)), nil
}
func (c Char) ToInt() (Int, error) {
	return Int(c), nil
}
func (c Char) ToPrint() string {
	return string(c)
}

/////////////////////////////
// String
func (b String) DataType() string { return DataTypeString }
func (s String) ToBool() (Bool, error) {
	switch strings.ToLower(string(s)) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("Cannot convert '%s' to Bool", s)
}
func (s String) ToChar() (Char, error) {
	return Char(s[0]), nil
}
func (s String) ToString() (String, error) {
	return s, nil
}
func (s String) ToInt() (Int, error) {
	n, err := strconv.Atoi(string(s))
	if err != nil {
		return 0, fmt.Errorf("Cannot convert '%s' to Int: %v", s, err)
	} else {
		return Int(n), nil
	}
}
func (s String) ToDouble() (Double, error) {
	var n float64
	_, err := fmt.Sscanf(string(s), "%f", &n)
	if err != nil {
		return 0, fmt.Errorf("Cannot convert '%s' to Double: %v", s, err)
	} else {
		return Double(n), nil
	}
}

func (s String) ToDateTime() (DateTime, error) {
	t, err := time.Parse("01/02/2006", string(s))
	return DateTime(t), err
}

/////////////////////////////
// Integer Range
func (n IntRange) DataType() string { return DataTypeIntRange }
func (n IntRange) ToString() (String, error) {
	return String(fmt.Sprintf("%d-%d", n.From, n.To)), nil
}

/////////////////////////////
// Conditional Values
func (n ConditionalValues) DataType() string { return DataTypeConditionalValues }
func (n ConditionalValues) ToString() (String, error) {
	ts, err := ToString(n.True)
	if err != nil {
		return "", err
	}
	fs, err := ToString(n.False)
	if err != nil {
		return "", err
	}
	return String(fmt.Sprintf("%s || %s", ts, fs)), nil
}
