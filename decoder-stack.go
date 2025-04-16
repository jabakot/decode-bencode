package decodebencode

import (
	"errors"
	"fmt"
	"slices"
)

type DataStack []interface{}

type dict_markerType struct{}
type list_markerType struct{}

// dictinoary makrer for stack
var DICT_MARKER = dict_markerType{}

// list makrer for stack
var LIST_MARKER = list_markerType{}

const (
	INT_CONTROL_SYMBOL   = 'i'
	STR_CONTROL_SYMBOL   = ':'
	LIST_CONTROL_SYMBOL  = 'l'
	DICT_CONTROL_SYMBOL  = 'd'
	CLOSE_CONTROL_SYMBOL = 'e'
)

func IsDictSymbol(cadidate interface{}) bool {
	_, ok := cadidate.(dict_markerType)
	return ok
}

func IsListSymbol(cadidate interface{}) bool {
	_, ok := cadidate.(list_markerType)
	return ok
}

func (s *DataStack) Push(v interface{}) {
	*s = append(*s, v)
}

func (s *DataStack) Pop() (interface{}, error) {
	l := len(*s)
	if l == 0 {
		return nil, errors.New("data stack is empty")
	}
	result := (*s)[l-1]
	*s = (*s)[:l-1]
	return result, nil
}

// Transforms sequence of elements in buff to map (bencode dictionary)
func ShrinkDictionary(stackSlice []interface{}) (map[string]interface{}, error) {
	if len(stackSlice)%2 != 0 {
		return nil, fmt.Errorf("cannot transform stackSlice to map because odd number of elements: %v (%d)", stackSlice, len(stackSlice))
	}
	slices.Reverse(stackSlice)
	dict := make(map[string]interface{})
	for i := 0; i < len(stackSlice); i += 2 {
		key, ok_key := stackSlice[i].(string)

		if !ok_key {
			return nil, fmt.Errorf("value %v cannot be used as dict key", stackSlice[i])
		}

		dict[key] = stackSlice[i+1]
	}

	return dict, nil
}

func ShrinkStack(stack *DataStack) error {
	buff := make([]interface{}, 0)
	stop := false

	for len(*stack) > 0 && !stop {
		el, err := stack.Pop()
		if err != nil {
			return err
		}

		if IsDictSymbol(el) {
			dict, err := ShrinkDictionary(buff)

			if err != nil {
				return err
			}
			*stack = append(*stack, dict)
			buff = make([]interface{}, 0)
			stop = true
		} else if IsListSymbol(el) {
			slices.Reverse(buff)
			*stack = append(*stack, buff)
			buff = make([]interface{}, 0)
			stop = true
		} else {
			buff = append(buff, el)
		}
	}

	if len(buff) == 1 {
		stack.Push(buff[0])
	}

	return nil
}
