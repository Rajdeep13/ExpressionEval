package evaluator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/contactkeval/expressioneval/datatype"
	"github.com/contactkeval/expressioneval/tokenizer"
)

// Supported intrinsic methods
var intrinsicMethods map[string]intrinsicMethod

// Set up all intrinsic methods
func init() {
	intrinsicMethods = map[string]intrinsicMethod{
		"Abs": polyTypeCheckedMethod(
			"N", func(args ...datatype.DataType) (datatype.DataType, error) {
				return datatype.Double(math.Abs(toFloat(args[0]))), nil
			}),
		"AddQuotes": polyTypeCheckedMethod(
			"S,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				s, b := toString(args[0]), toBool(args[1])
				if b {
					s = strings.Replace(s, "\\", "\\\\", -1)
					s = strings.Replace(s, "'", "\\'", -1)
					s = "'" + s + "'"
				} else {
					s = strconv.Quote(s)
				}
				return datatype.String(s), nil
			}),
		"Apy": polyTypeCheckedMethod(
			"N,N", func(args ...datatype.DataType) (datatype.DataType, error) {
				r, p := toFloat(args[0]), toFloat(args[1])
				return datatype.Double(100 * (math.Pow((1+(r/100)/p), p) - 1)), nil
			}),
		"Avg": polyTypeCheckedMethod(
			"LN", func(args ...datatype.DataType) (datatype.DataType, error) {
				l := toSlice(args[0])

				sum := 0.0
				for _, d := range l {
					sum += toFloat(d)
				}
				return datatype.Double(sum / float64(len(l))), nil
			}),
		"Contains": polyTypeCheckedMethod(
			"S,S,B", func(args ...datatype.DataType) (datatype.DataType, error) {
				haystack, needle, ignoreCase := toString(args[0]), toString(args[1]), toBool(args[2])
				if ignoreCase {
					haystack, needle = strings.ToLower(haystack), strings.ToLower(needle)
				}
				return datatype.Bool(strings.Contains(haystack, needle)), nil
			},
			"S,LS,BF,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				haystack, needles, ignoreCase, isAll := toString(args[0]), toSlice(args[1]), toBool(args[2]), toBool(args[3])
				if ignoreCase {
					haystack = strings.ToLower(haystack)
				}

				var contained int
				for _, needleVal := range needles {
					needle := toString(needleVal)
					if ignoreCase {
						needle = strings.ToLower(needle)
					}
					if strings.Contains(haystack, needle) {
						contained++
					}
				}

				if isAll {
					return datatype.Bool(contained == len(needles)), nil
				}
				return datatype.Bool(contained > 0), nil
			},
			"L,L,BF,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				haystacks, needles, ignoreCase, isAll := toSlice(args[0]), toSlice(args[1]), toBool(args[2]), toBool(args[3])

				var contained int
				for _, haystack := range haystacks {
					if s, ok := haystack.(datatype.String); ok && ignoreCase {
						haystack = datatype.String(strings.ToLower(toString(s)))
					}

					for _, needle := range needles {
						if s, ok := needle.(datatype.String); ok && ignoreCase {
							needle = datatype.String(strings.ToLower(toString(s)))
						}
						if reflect.DeepEqual(haystack, needle) {
							contained++
							break
						}
					}
				}

				if isAll {
					return datatype.Bool(contained == len(needles)), nil
				}
				return datatype.Bool(contained > 0), nil
			},
		),
		"Distinct": polyTypeCheckedMethod(
			"L", func(args ...datatype.DataType) (datatype.DataType, error) {
				l := toSlice(args[0])
				var seen []datatype.DataType

			Outer:
				for _, item := range l {
					for _, s := range seen {
						if reflect.DeepEqual(item, s) {
							continue Outer
						}
					}
					seen = append(seen, item)
				}

				return datatype.List(seen), nil
			}),
		"Dpr": polyTypeCheckedMethod(
			"N,I0", func(args ...datatype.DataType) (datatype.DataType, error) {
				r, y := toFloat(args[0]), toInt(args[1])
				if y == 0 {
					y = time.Now().Year()
				}
				t, _ := time.Parse("2006-01-02", strconv.Itoa(y)+"-12-31")

				return datatype.Double(r / (float64(t.YearDay()) * 100)), nil
			}),
		"EndsWith": polyTypeCheckedMethod(
			"S,S,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, part, ignoreCase := toString(args[0]), toString(args[1]), toBool(args[2])
				if ignoreCase {
					str, part = strings.ToLower(str), strings.ToLower(part)
				}
				return datatype.Bool(strings.HasSuffix(str, part)), nil
			},
			"S,L,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, parts, ignoreCase := toString(args[0]), toSlice(args[1]), toBool(args[2])
				if ignoreCase {
					str = strings.ToLower(str)
				}
				for _, partVal := range parts {
					part := toString(partVal)
					if ignoreCase {
						part = strings.ToLower(part)
					}
					if strings.HasSuffix(str, part) {
						return datatype.Bool(true), nil
					}
				}
				return datatype.Bool(false), nil
			},
		),
		"In": polyTypeCheckedMethod(
			"S,LS,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, l, ignoreCase := toString(args[0]), toSlice(args[1]), toBool(args[2])
				if ignoreCase {
					str = strings.ToLower(str)
				}

				for _, item := range l {
					itemStr := toString(item)
					if ignoreCase {
						itemStr = strings.ToLower(toString(item))
					}
					if reflect.DeepEqual(str, itemStr) {
						return datatype.Bool(true), nil
					}
				}

				return datatype.Bool(false), nil
			},
			"L,L,BF,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				haystacks, needles, ignoreCase, isAll := toSlice(args[0]), toSlice(args[1]), toBool(args[2]), toBool(args[3])

				var contained int
				for _, haystack := range haystacks {
					if s, ok := haystack.(datatype.String); ok && ignoreCase {
						haystack = datatype.String(strings.ToLower(toString(s)))
					}

					for _, needle := range needles {
						if s, ok := needle.(datatype.String); ok && ignoreCase {
							needle = datatype.String(strings.ToLower(toString(s)))
						}
						if reflect.DeepEqual(haystack, needle) {
							contained++
							break
						}
					}
				}

				if isAll {
					return datatype.Bool(contained == len(needles)), nil
				}
				return datatype.Bool(contained > 0), nil
			},
		),
		"IndexOf": polyTypeCheckedMethod(
			"S,S,I0,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, part, startIndex, ignoreCase := toString(args[0]), toString(args[1]), toInt(args[2]), toBool(args[3])
				if ignoreCase {
					str, part = strings.ToLower(str), strings.ToLower(part)
				}

				idx := strings.Index(str[startIndex:], part)
				if idx >= 0 {
					idx += startIndex
				}
				return datatype.Int(idx), nil
			},
			"S,LS,I0,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, parts, startIndex, ignoreCase := toString(args[0]), toSlice(args[1]), toInt(args[2]), toBool(args[3])
				if ignoreCase {
					str = strings.ToLower(str)
				}

				var idx = len(str)

				for _, partVal := range parts {
					part := toString(partVal)
					if ignoreCase {
						part = strings.ToLower(part)
					}
					if partIdx := strings.Index(str[startIndex:], part); partIdx >= 0 {
						if partIdx < idx {
							idx = partIdx
						}
					}
				}

				if idx == len(str) {
					idx = -1
				} else {
					idx += startIndex
				}
				return datatype.Int(idx), nil
			},
		),
		"Interval": polyTypeCheckedMethod(
			"", func(args ...datatype.DataType) (datatype.DataType, error) {
				return nil, fmt.Errorf("Not implemented")
			}),
		"JsonSelect": polyTypeCheckedMethod(
			"S,S", func(args ...datatype.DataType) (datatype.DataType, error) {
				j, sel := toString(args[0]), toString(args[1])

				var jm map[string]interface{}
				if err := json.Unmarshal([]byte(j), &jm); err != nil {
					return nil, err
				}

				parts := strings.Split(strings.Replace(sel, "[", ".[", -1), ".")

				var node interface{} = jm

				for _, part := range parts {
					if len(part) == 0 {
						continue
					}
					if part[0] == '[' {
						if arr, ok := node.([]interface{}); ok {
							idxStr := strings.Replace(strings.Replace(part, "[", "", -1), "]", "", -1)

							if idx, err := strconv.Atoi(idxStr); err != nil {
								return nil, fmt.Errorf("'%s' is not a valid index", idx)
							} else {
								node = arr[idx]
							}
						} else {
							return nil, fmt.Errorf("Using array index for non-array node: %v", node)
						}
					} else {
						if m, ok := node.(map[string]interface{}); ok {
							node = m[part]
						} else {
							return nil, fmt.Errorf("Using map index for non-map node: %v", node)
						}
					}
				}

				nodeJSON, _ := json.Marshal(node)

				return datatype.String(string(nodeJSON)), nil
			}),
		"GetWebPage": polyTypeCheckedMethod(
			"S,S", func(args ...datatype.DataType) (datatype.DataType, error) {
				url, except := toString(args[0]), toString(args[1])
				if resp, err := http.Get(url); err != nil {
					return datatype.String(except), nil
				} else {
					if body, err := ioutil.ReadAll(resp.Body); err != nil {
						return datatype.String(except), nil
					} else {
						return datatype.String(string(body)), nil
					}
				}
			}),
		"Length": polyTypeCheckedMethod(
			"S", func(args ...datatype.DataType) (datatype.DataType, error) {
				str := toString(args[0])
				return datatype.Int(len(str)), nil
			},
			"L", func(args ...datatype.DataType) (datatype.DataType, error) {
				l := toSlice(args[0])
				return datatype.Int(len(l)), nil
			},
		),
		"Matches": polyTypeCheckedMethod(
			"S,S", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, patStr := toString(args[0]), toString(args[1])

				if matched, err := regexp.MatchString(str, patStr); err != nil {
					return nil, err
				} else {
					return datatype.Bool(matched), nil
				}
			}),
		"Max": polyTypeCheckedMethod(
			"LN", func(args ...datatype.DataType) (datatype.DataType, error) {
				l := toSlice(args[0])

				max := toFloat(l[0])
				for _, d := range l {
					if toFloat(d) > max {
						max = toFloat(d)
					}
				}
				return datatype.Double(max), nil
			}),
		"Med": polyTypeCheckedMethod(
			"LN", func(args ...datatype.DataType) (datatype.DataType, error) {
				l := toSlice(args[0])
				var ln []float64
				for _, item := range l {
					ln = append(ln, toFloat(item))
				}
				sort.Float64s(ln)
				if len(ln)%2 == 0 {
					return datatype.Double((ln[len(ln)/2-1] + ln[len(ln)/2]) / 2), nil
				} else {
					return datatype.Double(ln[len(ln)/2]), nil
				}
			}),
		"Min": polyTypeCheckedMethod(
			"LN", func(args ...datatype.DataType) (datatype.DataType, error) {
				l := toSlice(args[0])

				min := toFloat(l[0])
				for _, d := range l {
					if toFloat(d) < min {
						min = toFloat(d)
					}
				}
				return datatype.Double(min), nil
			}),
		"Piece": polyTypeCheckedMethod(
			"S,S,I0,I0", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, delim, startCount, lastCount := toString(args[0]), toString(args[1]), toInt(args[2]), toInt(args[3])
				parts := strings.Split(str, delim)

				if lastCount == 0 {
					lastCount = len(parts)
				}
				startCount--
				lastCount--

				if len(parts) <= startCount {
					return datatype.String(""), nil
				}
				if len(parts) <= lastCount {
					lastCount = len(parts) - 1
				}

				return datatype.String(strings.Join(parts[startCount:lastCount+1], delim)), nil
			}),
		"Pow": polyTypeCheckedMethod(
			"N,N", func(args ...datatype.DataType) (datatype.DataType, error) {
				n, p := toFloat(args[0]), toFloat(args[1])
				return datatype.Double(math.Pow(n, p)), nil
			}),
		"Pv": polyTypeCheckedMethod(
			"N,N,I1", func(args ...datatype.DataType) (datatype.DataType, error) {
				n, r, p := toFloat(args[0]), toFloat(args[1]), toInt(args[2])
				return datatype.Double(n / math.Pow(1+(r/100), float64(p))), nil
			}),
		"Pva": polyTypeCheckedMethod(
			"N,N,I1", func(args ...datatype.DataType) (datatype.DataType, error) {
				n, r, p := toFloat(args[0]), toFloat(args[1]), toInt(args[2])
				return datatype.Double(n * ((1 - math.Pow(1+(r/100), float64(-p))) / (r / 100))), nil
			}),
		"Replace": polyTypeCheckedMethod(
			"S,S,S", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, from, to := toString(args[0]), toString(args[1]), toString(args[2])
				return datatype.String(strings.Replace(str, from, to, -1)), nil
			}),
		"ShowTokens": polyTypeCheckedMethod(
			"S", func(args ...datatype.DataType) (datatype.DataType, error) {
				str := toString(args[0])
				tok, err := tokenizer.Tokenize(str)
				return datatype.String(fmt.Sprintf("%#v", tok)), err
			}),
		"Split": polyTypeCheckedMethod(
			"S,C", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, ch := toString(args[0]), toRune(args[1])
				parts := strings.Split(str, string(ch))
				l := datatype.List{}
				for _, part := range parts {
					l = append(l, datatype.String(part))
				}
				return l, nil
			}),
		"Sort": polyTypeCheckedMethod(
			"L", func(args ...datatype.DataType) (datatype.DataType, error) {
				l := toSlice(args[0])
				res := datatype.List{}

				switch t := l[0].DataType(); t {
				case datatype.DataTypeString:
					ls := toStringSlice(l)
					sort.Strings(ls)
					for _, item := range ls {
						res = append(res, datatype.String(item))
					}
				case datatype.DataTypeChar:
					ls := toStringSlice(l)
					sort.Strings(ls)
					for _, item := range ls {
						res = append(res, datatype.Char(item[0]))
					}
				case datatype.DataTypeInt:
					ls := toIntSlice(l)
					sort.Ints(ls)
					for _, item := range ls {
						res = append(res, datatype.Int(item))
					}
				case datatype.DataTypeDouble:
					ls := toFloatSlice(l)
					sort.Float64s(ls)
					for _, item := range ls {
						res = append(res, datatype.Double(item))
					}
				default:
					return nil, fmt.Errorf("Cannot sort  data type %s", t)
				}
				return res, nil
			}),
		"StartsWith": polyTypeCheckedMethod(
			"S,S,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, part, ignoreCase := toString(args[0]), toString(args[1]), toBool(args[2])
				if ignoreCase {
					str, part = strings.ToLower(str), strings.ToLower(part)
				}
				return datatype.Bool(strings.HasPrefix(str, part)), nil
			},
			"S,L,BF", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, parts, ignoreCase := toString(args[0]), toSlice(args[1]), toBool(args[2])
				if ignoreCase {
					str = strings.ToLower(str)
				}
				for _, partVal := range parts {
					part := toString(partVal)
					if ignoreCase {
						part = strings.ToLower(part)
					}
					if strings.HasPrefix(str, part) {
						return datatype.Bool(true), nil
					}
				}
				return datatype.Bool(false), nil
			},
		),
		"StripQuotes": polyTypeCheckedMethod(
			"S", func(args ...datatype.DataType) (datatype.DataType, error) {
				str := toString(args[0])
				ustr, err := strconv.Unquote(str)
				return datatype.String(ustr), err
			}),
		"ToString": intrinsicMethodFunc(func(pfe *PostfixExpression) (datatype.DataType, error) {
			funcArg, err := GetUnaryOperand(pfe)
			if err != nil {
				return nil, err
			}

			if vs, ok := funcArg.(interface {
				ToString() (datatype.String, error)
			}); ok {
				s, err := vs.ToString()
				return s, err
			}

			return nil, fmt.Errorf("Could not convert %v to string", funcArg.DataType())
		}),
		"ToUpper": polyTypeCheckedMethod(
			"S", func(args ...datatype.DataType) (datatype.DataType, error) {
				str := toString(args[0])
				ustr := strings.ToUpper(str)
				return datatype.String(ustr), nil
			}),
		"Translate": polyTypeCheckedMethod(
			"S,S,S", func(args ...datatype.DataType) (datatype.DataType, error) {
				str, from, to := toString(args[0]), toString(args[1]), toString(args[2])
				if len(from) != len(to) {
					return nil, errors.New("'From' and 'To' arguments should have equal lengths")
				}
				m := map[rune]rune{}
				var frunes, trunes []rune
				for _, c := range from {
					frunes = append(frunes, c)
				}
				for _, c := range to {
					trunes = append(trunes, c)
				}
				for i, _ := range from {
					m[frunes[i]] = trunes[i]
				}

				var tstr string
				for _, ch := range str {
					if t, ok := m[ch]; ok {
						tstr += string(t)
					} else {
						tstr += string(ch)
					}
				}

				return datatype.String(tstr), nil
			}),
	}
}
