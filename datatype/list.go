package datatype

import "fmt"

// A list of values. The type of the list is determined by the type of the
// first element in the list.
type List []DataType

func (l List) DataType() string {
	if len(l) > 0 {
		return fmt.Sprintf("%s[]", l[0].DataType())
	} else {
		return "Unknown[]"
	}
}

type converter func(DataType) (DataType, error)

// convert all elements of a list to a particular type using the convert
// function
func (l List) convertDataType(convert converter) (List, error) {
	nl := List{}
	for _, item := range l {
		v, err := convert(item)
		if err != nil {
			return nil, err
		}
		nl = append(nl, v)
	}
	return nl, nil
}

func (l List) AllToBool() (List, error) {
	return l.convertDataType(func(v DataType) (DataType, error) { cv, err := ToBool(v); return cv, err })
}

func (l List) AllToInt() (List, error) {
	return l.convertDataType(func(v DataType) (DataType, error) { cv, err := ToInt(v); return cv, err })
}

func (l List) AllToDouble() (List, error) {
	return l.convertDataType(func(v DataType) (DataType, error) { cv, err := ToDouble(v); return cv, err })
}

func (l List) AllToString() (List, error) {
	return l.convertDataType(func(v DataType) (DataType, error) { cv, err := ToString(v); return cv, err })
}

func (l List) AllToChar() (List, error) {
	return l.convertDataType(func(v DataType) (DataType, error) { cv, err := ToChar(v); return cv, err })
}

func (l List) ToString() (String, error) {
	res := ""
	for i, item := range l {
		if i > 0 {
			res += ", "
		}
		res += ToPrint(item)
	}

	return String(res), nil
}

func (l List) ToPrint() string {
	res := "["

	for i, item := range l {
		if i > 0 {
			res += ", "
		}
		res += ToPrint(item)
	}
	res += "]"

	return res
}

// CommaList is a list of operands specified using the comma operator.
// It is mainly used when evaluating tokens as the output of the comma
// operator. A CommaList can contain CommaLists, but they're only useful during
// parsing. Call the Flatten() method to convert it into a List before
// performing any operations.
type CommaList []DataType

func (cl CommaList) DataType() string {
	return "comma_list"
}

func (cl CommaList) Flatten() List {
	l := List{}

	for _, item := range cl {
		switch v := item.(type) {
		case CommaList:
			l2 := v.Flatten()
			l = append(l, l2...)
		default:
			l = append(l, v)
		}
	}

	return l
}
