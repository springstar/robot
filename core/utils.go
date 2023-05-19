package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func SetField(obj interface{}, name string, value interface{}) error {

    structValue := reflect.ValueOf(obj).Elem()
    fieldVal := structValue.FieldByName(name)

    if !fieldVal.IsValid() {
        return fmt.Errorf("No such field: %s in obj", name)
    }

    if !fieldVal.CanSet() {
        return fmt.Errorf("Cannot set %s field value", name)
    }

    val := reflect.ValueOf(value)

    if fieldVal.Type() != val.Type() {

        if m,ok := value.(map[string]interface{}); ok {

            // if field value is struct
            if fieldVal.Kind() == reflect.Struct {
                return FillStruct(m, fieldVal.Addr().Interface())
            }

            // if field value is a pointer to struct
            if fieldVal.Kind()==reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
                if fieldVal.IsNil() {
                    fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
                }
                // fmt.Printf("recursive: %v %v\n", m,fieldVal.Interface())
                return FillStruct(m, fieldVal.Interface())
            }

        }

        return fmt.Errorf("Provided value type didn't match obj field type")
    }

    fieldVal.Set(val)
    return nil

}

func FillStruct(m map[string]interface{}, s interface{}) error {
    for k, v := range m {
        err := SetField(s, k, v)
        if err != nil {
            return err
        }
    }
    return nil
}


func GetType(v interface{}) string {
    if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
        return "*" + t.Elem().Name()
    } else {
        return t.Name()
    }
}

func Str2Int(s string) (int, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return -1, err
	}

	return int(n), nil
}

func Str2IntSlice(s string) ([]int, error) {
	strSlice := strings.Split(s, ",")
	intSlice := make([]int, len(strSlice))
	for _, v := range strSlice {
		n, err := Str2Int(v)
		if err != nil {
			return nil, err
		}

		intSlice = append(intSlice, int(n))
	}

	return intSlice, nil
}

func Str2Int32Slice(s string) ([]int32, error) {
	strSlice := strings.Split(s, ",")
	intSlice := make([]int32, len(strSlice))
	for _, v := range strSlice {
		n, err := Str2Int(v)
		if err != nil {
			return nil, err
		}

		intSlice = append(intSlice, int32(n))
	}

	return intSlice, nil
}