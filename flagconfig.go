package envconfig

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	// "github.com/davidwalter0/envconfig/flag"
	"./flag"
)

// Parse calls ParseEnv then ParseFlags
func Parse(prefix string, T interface{}) error {
	return Process(prefix, T)
}

// Process wraps Parse which orders sets precedence
func Process(prefix string, T interface{}) error {
	err := ParseEnv(prefix, T)
	err = ParseFlags(T)
	flag.Parse()
	return err
}

// ParseFlags wraps BuildFlagsFromStructPtr
func ParseFlags(T interface{}) error {
	return BuildFlagsFromStructPtr(T)
}

// TypeAlignDefault for default setup returning interface and error
func TypeAlignDefault(text string, T reflect.StructField) (interface{}, error) {
	switch T.Type.Kind() {
	case reflect.String:
		return text, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(text) == 0 {
			return 0, nil
		}
		if T.Type.Kind() == reflect.Int64 &&
			T.Type.PkgPath() == "time" &&
			T.Type.Name() == "Duration" {
			return time.ParseDuration(text)
		}
		return strconv.ParseInt(text, 0, T.Type.Bits())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.ParseUint(text, 0, T.Type.Bits())
	case reflect.Bool:
		return strconv.ParseBool(text)
	case reflect.Float32, reflect.Float64:
		lhs, rhs := strconv.ParseFloat(text, T.Type.Bits())
		return (float64)(lhs), rhs
	case reflect.Slice:
		fmt.Println("fix map in next iteration", text)
		lhs := flag.SliceValue(strings.Split(text, ","))
		fmt.Println("fix map in next iteration", lhs)
		return lhs, nil
	case reflect.Map:
		fmt.Println("fix map in next iteration")
	}
	return nil, nil

}

// HyphenateCamelCase converts camel case name string and hyphenates
// words for flags between words
func HyphenateCamelCase(name string) (key string) {
	expr := regexp.MustCompile("([^A-Z]+|[A-Z][^A-Z]+|[A-Z]+)")

	words := expr.FindAllStringSubmatch(name, -1)
	if len(words) > 0 {
		var name []string
		for _, words := range words {
			name = append(name, strings.ToLower(words[0]))
		}
		key = strings.Join(name, "-")
	}
	return
}

// BuildFlagsFromStructPtr takes a pointer to struct and generates
// flags using the struct tags as supplemental meta data.
// To support a prioritized ordered assignment
// 1. flag setting      : overrides
// 2. struct default tag: overrides
// 3. Environment variables
//
// ProcessEnv is called before flags, if a value is set from the
// environment
//
// Assign default value if declared in the StructTag and create
// the new flag, create short name and long name if tags exist
//
// default to the lowercase of the variable name, Abc to -abc
//
// e.g. CamelCase use downcase hyphenated breaks for the name:
// -camel-case
func BuildFlagsFromStructPtr(T interface{}) error {
	if reflect.ValueOf(T).Kind() != reflect.Ptr {
		log.Fatalln("Unable to build flags from a struct, pass a *struct ptr")
	}
	// Dereference into an adressable value
	element := reflect.ValueOf(T).Elem()
	elementType := element.Type()
	for i := 0; i < elementType.NumField(); i++ {
		field := elementType.Field(i)
		// Get tags for this field
		name := field.Tag.Get("name")
		short := field.Tag.Get("short")
		usage := field.Tag.Get("usage")
		text := field.Tag.Get("default")
		parsedName := HyphenateCamelCase(field.Name)

		addr := element.Field(i).Addr().Interface()
		if len(name) == 0 {
			if len(parsedName) == 0 {
				log.Fatalln("Name failed for target")
			}
			name = parsedName
		}

		var defaultValue interface{}
		var err error
		if len(text) != 0 {
			defaultValue, err = TypeAlignDefault(text, field)
			if err != nil {
				log.Fatalln("Error", err)
			}
		}

		names := []string{}
		if len(name) > 0 {
			names = append(names, name)
		}
		if len(short) > 0 {
			names = append(names, short)
		}

		for _, arg := range names {
			switch ptr := addr.(type) {
			case *time.Duration:
				var v time.Duration
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(time.Duration)
					}
				}
				flag.DurationVar(ptr, arg, v, usage)
			case *bool:
				var v bool
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(bool)
					}
				}
				flag.BoolVar(ptr, arg, v, usage)
			case *float32:
				var v float32
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(float32)
					}
				}
				flag.Float32Var(ptr, arg, v, usage)
			case *float64:
				var v float64
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(float64)
					}
				}
				flag.Float64Var(ptr, arg, v, usage)
			case *int:
				var v int
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = (int)(defaultValue.(int64))
					}
				}
				flag.IntVar(ptr, arg, v, usage)
			case *uint:
				var v uint
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = (uint)(defaultValue.(uint64))
					}
				}
				flag.UintVar(ptr, arg, v, usage)
			case *int64:
				var v int64
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(int64)
					}
				}
				flag.Int64Var(ptr, arg, v, usage)
			case *uint64:
				var v uint64
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(uint64)
					}
				}
				flag.Uint64Var(ptr, arg, v, usage)
			case *int8:
				var v int8
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(int8)
					}
				}
				flag.Int8Var(ptr, arg, v, usage)
			case *uint8:
				var v uint8
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(uint8)
					}
				}
				flag.Uint8Var(ptr, arg, v, usage)
			case *int16:
				var v int16
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(int16)
					}
				}
				flag.Int16Var(ptr, arg, v, usage)
			case *uint16:
				var v uint16
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(uint16)
					}
				}
				flag.Uint16Var(ptr, arg, v, usage)
			case *int32:
				var v int32
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(int32)
					}
				}
				flag.Int32Var(ptr, arg, v, usage)
			case *uint32:
				var v uint32
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(uint32)
					}
				}
				flag.Uint32Var(ptr, arg, v, usage)
			case *[]string:
				var v flag.SliceValue
				if defaultValue != nil {
					sliceValue := defaultValue.(flag.SliceValue)
					if len(sliceValue) > 0 {
						for _, u := range sliceValue {
							v = append(v, u)
						}
					}
					v = defaultValue.(flag.SliceValue)
				}
				flag.SliceVar((*flag.SliceValue)(ptr), arg, v, usage)
			case *string:
				var v string
				if *ptr != v {
					v = *ptr
				} else {
					if defaultValue != nil {
						v = defaultValue.(string)
					}
				}
				flag.StringVar(ptr, arg, v, usage)
			default:
				log.Printf("unknown/default type:%v %T %p\n", ptr, ptr, ptr)
			}
		}
	}
	return nil
}
