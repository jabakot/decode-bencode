package decodebencode

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
)

// yes, length is in bytes
func EncodeBencodeString(s string) string {
	if len(s) < 1 {
		return ""
	}
	return fmt.Sprintf("%d:%s", len(s), s)
}

func EncodeBencodeInteger(i int) string {
	return fmt.Sprintf("i%de", i)
}

func EncodeBencodeList(list []any) string {

	if len(list) == 1 {
		return "le"
	}

	buff := "l"
	for _, v := range list {
		t := reflect.TypeOf(v)
		kind := t.Kind()

		switch kind {
		case reflect.Array, reflect.Slice:
			arr_val, ok := v.([]any)
			if ok {
				buff += EncodeBencodeList(arr_val)
			} else {

			}
		case reflect.Int:
			buff += EncodeBencodeInteger(int(reflect.ValueOf(v).Int()))
		case reflect.String:
			buff += EncodeBencodeString(reflect.ValueOf(v).String())
		case reflect.Map:
			dict_val, ok := v.(map[string]any)
			if ok {
				buff += EncodeBencodeDict(dict_val)
			} else {
				fmt.Println("cannot cast map", v)
			}

		default:
			fmt.Printf("cannot process element:  %v", v)
		}

	}
	buff += "e"

	return buff
}

func EncodeBencodeDict(dict map[string]any) string {
	keys := slices.Sorted(maps.Keys(dict))

	if len(keys) == 0 {
		return "de"
	}

	buff := "d"

	for _, key := range keys {
		val := dict[key]
		t := reflect.TypeOf(dict[key])
		kind := t.Kind()
		switch kind {
		case reflect.Array, reflect.Slice:
			arr_val, ok := val.([]any)
			if ok {
				buff += EncodeBencodeString(key)
				buff += EncodeBencodeList(arr_val)
			} else {
				fmt.Println("cannot cast array or slice ", val)
			}
		case reflect.Int:
			buff += EncodeBencodeString(key)
			buff += EncodeBencodeInteger(int(reflect.ValueOf(val).Int()))
		case reflect.String:
			buff += EncodeBencodeString(key)
			buff += EncodeBencodeString(reflect.ValueOf(val).String())
		case reflect.Map:
			dict_val, ok := val.(map[string]any)
			if ok {
				buff += EncodeBencodeString(key)
				buff += EncodeBencodeDict(dict_val)
			} else {
				fmt.Println("cannot cast map ", val)
			}
		default:
			fmt.Printf("cannot process element: [%v] : %v", key, val)
		}

	}

	buff += "e"

	return buff
}
