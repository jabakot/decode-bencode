package decodebencode

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// finds next [r] rune int input, useful for
func FindNextRune(start int, r rune, input []rune) (int, bool) {
	inputLength := len(input)

	if start >= inputLength {
		return -1, false
	}

	for i := start; i < inputLength; i++ {
		if input[i] == r {
			return i, true
		}
	}

	return -1, false
}

func ParseInt(input []rune) (int, error) {
	str := string(input)
	num, err := strconv.Atoi(str)

	if err != nil {
		return -1, errors.New("cannot convert to int")
	}

	return num, nil
}

func DecodeBencode(input string) (interface{}, error) {
	if len(strings.TrimSpace(input)) == 0 {
		return nil, nil
	}

	i := 0
	r_input := []rune(input)
	stack := make(DataStack, 0, 1)

	for i < len(r_input) {
		switch r_input[i] {
		case INT_CONTROL_SYMBOL:
			i++
			start := i

			end, found_end := FindNextRune(start, CLOSE_CONTROL_SYMBOL, r_input)

			if !found_end {
				return nil, fmt.Errorf("cannot find closing symbol %v for integer, starting from: %d in %s", CLOSE_CONTROL_SYMBOL, start, string(r_input))
			}

			num, err_parse_int := ParseInt(r_input[start:end])

			if err_parse_int != nil {
				fmt.Println(err_parse_int)
				return nil, err_parse_int
			}

			stack.Push(num)
			i = end + 1

		case LIST_CONTROL_SYMBOL:
			stack.Push(LIST_MARKER)
			i++

		case DICT_CONTROL_SYMBOL:
			stack.Push(DICT_MARKER)
			i++

		case CLOSE_CONTROL_SYMBOL:
			if i < len(r_input) {
				zip_error := ShrinkStack(&stack)

				if zip_error != nil {
					fmt.Println(zip_error)
					return nil, zip_error
				}
			}
			i++
			// try to parse string
		default:
			if !unicode.IsDigit(r_input[i]) {
				return nil, fmt.Errorf("parsing error, expected digit, got %v on index %d", string(r_input[i-2:]), i)
			}

			start := i
			semicolon_index, found_semicolon_index := FindNextRune(start, STR_CONTROL_SYMBOL, r_input)

			if !found_semicolon_index {
				return nil, fmt.Errorf("cannot find closing symbol %v for string, starting from: %d in %s", STR_CONTROL_SYMBOL, start, string(r_input))
			}

			str_bytes_length, str_bytes_length_error := ParseInt(r_input[start:semicolon_index])
			if str_bytes_length_error != nil {
				fmt.Println(string(r_input[start:semicolon_index]))
				return nil, str_bytes_length_error
			}

			if len(string(r_input[semicolon_index+1:])) < str_bytes_length {
				return nil, fmt.Errorf("wrong string encoding: length of string %d is greater than remainng length of %v", str_bytes_length, string(r_input[semicolon_index+1:]))
			}

			str_start_index := semicolon_index + 1

			j := 0
			str_end_index := str_start_index
			for j < str_bytes_length {
				j += utf8.RuneLen(r_input[str_end_index])
				str_end_index++
			}

			str := string(r_input[str_start_index:str_end_index])

			stack.Push(str)
			i = str_end_index
		}
	}

	if len(stack) > 1 {
		return nil, fmt.Errorf("wrong input data, faced sequence of unwrapped elements: %v", stack)
	}

	el, err := stack.Pop()

	if err != nil {
		return nil, err
	}

	return el, nil
}
