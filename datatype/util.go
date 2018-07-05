package datatype

import "fmt"

func IsDateTime(vs ...DataType) bool {
	for _, v := range vs {
		if v.DataType() != DataTypeDateTime {
			return false
		}
	}
	return true
}

func IsInt(vs ...DataType) bool {
	for _, v := range vs {
		if v.DataType() != DataTypeInt {
			return false
		}
	}
	return true
}

func IsNumber(vs ...DataType) bool {
	for _, v := range vs {
		if v.DataType() != DataTypeInt && v.DataType() != DataTypeDouble {
			return false
		}
	}
	return true
}

func IsString(vs ...DataType) bool {
	for _, v := range vs {
		if v.DataType() != DataTypeString {
			return false
		}
	}
	return true
}

func IsBool(vs ...DataType) bool {
	for _, v := range vs {
		if v.DataType() != DataTypeBool {
			return false
		}
	}
	return true
}

func IsListAny(vs ...DataType) bool {
	res := false
	for _, v := range vs {
		_, isList := v.(List)
		res = res || isList
	}
	return res
}

func IsListAll(vs ...DataType) bool {
	for _, v := range vs {
		if _, isList := v.(List); !isList {
			return false
		}
	}
	return true
}

// Type upgrade helpers

// not really an upgrade
func UpgradeBinaryOperandsToStrings(op1, op2 DataType) (String, String, error) {
	if !IsString(op1, op2) {
		return "", "", fmt.Errorf("Cannot upgrade to String: Not all operands are strings")
	}
	sop1, err := ToString(op1)
	if err != nil {
		return "", "", err
	}
	sop2, err := ToString(op2)
	return sop1, sop2, err
}

func UpgradeBinaryOperandsToDoubles(op1, op2 DataType) (Double, Double, error) {
	if !IsNumber(op1, op2) {
		return 0, 0, fmt.Errorf("Cannot upgrade to Double: Not all operands are numbers")
	}
	dop1, err := ToDouble(op1)
	if err != nil {
		return 0, 0, err
	}
	dop2, err := ToDouble(op2)
	return dop1, dop2, err
}

func UpgradeBinaryOperandsToLists(op1, op2 DataType) (List, List, error) {
	if !IsListAny(op1, op2) {
		return nil, nil, fmt.Errorf("Cannot upgrade to List: At least one operand should be a list")
	}

	var lop1, lop2 List

	if v, ok := op1.(List); ok {
		lop1 = v
	} else {
		lop1 = List{op1}
	}

	if v, ok := op2.(List); ok {
		lop2 = v
	} else {
		lop2 = List{op2}
	}

	return lop1, lop2, nil
}

// Helper method to convert a data type into a printable string
func ToPrint(op DataType) string {
	switch v := op.(type) {
	case Printable:
		return v.ToPrint()
	default:
		return fmt.Sprintf("%v", v)
	}
}
