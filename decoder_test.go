package decodebencode_test

import (
	"reflect"
	"testing"

	decodebencode "github.com/jabakot/decode-bencode"
)

func TestFindNextRune(t *testing.T) {
	runes := []rune("hello")

	type TestCase struct {
		name       string
		start      int
		targetRune rune
		input      []rune
		expected   int
		ok         bool
	}

	testCases := []TestCase{
		{name: "Finds `e` rune in hello", start: 0, targetRune: 'e', input: runes, expected: 1, ok: true},
		{name: "Fails to find `ゴ` rune in hello", start: 0, targetRune: 'ゴ', input: runes, expected: -1, ok: false},
		{name: "Start is out of input length", start: 10, targetRune: 'e', input: runes, expected: -1, ok: false},
		{name: "Empty input", start: 0, targetRune: 'a', input: []rune{}, expected: -1, ok: false},
		{name: "Multiple matches, return first from start", start: 0, targetRune: 'a', input: []rune("tatata"), expected: 1, ok: true},
		{name: "Match is exactly at start", start: 0, targetRune: 'h', input: runes, expected: 0, ok: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			output, ok := decodebencode.FindNextRune(tc.start, tc.targetRune, tc.input)
			if output != tc.expected {
				t.Errorf("Expected value: %d,got: %d", tc.expected, output)
			}
			if ok != tc.ok {
				t.Errorf("Expected status: %v,got: %v", tc.ok, ok)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	type TestCase struct {
		name      string
		input     []rune
		expected  int
		expectErr bool
	}

	testCases := []TestCase{
		{
			name:     "Valid positive number",
			input:    []rune("123"),
			expected: 123,
		},
		{
			name:     "Valid negative number",
			input:    []rune("-456"),
			expected: -456,
		},
		{
			name:      "Invalid input: letters",
			input:     []rune("abc"),
			expectErr: true,
		},
		{
			name:      "Invalid input: mixed letters and numbers",
			input:     []rune("12ab"),
			expectErr: true,
		},
		{
			name:      "Empty input",
			input:     []rune(""),
			expectErr: true,
		},
		{
			name:     "Zero",
			input:    []rune("0"),
			expected: 0,
		},
		{
			name:      "Leading whitespace",
			input:     []rune(" 42"),
			expectErr: true, // strconv.Atoi doesn't allow leading spaces
		},
		{
			name:      "Trailing newline",
			input:     []rune("42\n"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := decodebencode.ParseInt(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil, result: %d", output)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if output != tc.expected {
					t.Errorf("Expected %d, got %d", tc.expected, output)
				}
			}
		})
	}

}

func TestDecodeBencode(t *testing.T) {
	type TestCase struct {
		name      string
		input     string
		expected  interface{}
		expectErr bool
	}

	test_map_simple := map[string]interface{}{"answer": 42, "hello": "world", "hi": "mark", "jojo": "ゴゴゴゴ", "wrong-answer": -42}
	test_list_simple := []interface{}{42, "hi"}

	test_list_complex := []interface{}{42, "hi", test_map_simple, test_list_simple}
	test_map_complex := map[string]interface{}{"answer": 42, "hello": "world", "list": test_list_simple, "dict": test_map_simple}

	testCases := []TestCase{
		{name: "empty input", input: "", expected: nil, expectErr: false},
		{name: "empty input with spaces", input: "                      ", expected: nil, expectErr: false},
		{name: "integer", input: "i42e", expected: 42, expectErr: false},
		{name: "integer without close symbol", input: "i42", expected: nil, expectErr: true},
		{name: "integer with wrong data inside", input: "ihi!e", expected: nil, expectErr: true},
		{name: "string", input: "2:hi", expected: "hi", expectErr: false},
		{name: "string with incorrect byte length", input: "4:hi", expected: "hi", expectErr: true},
		{name: "string with incorrect byte length", input: "1:hi", expected: "hi", expectErr: true},
		{name: "unicode string", input: "12:ゴゴゴゴ", expected: "ゴゴゴゴ", expectErr: false},
		{name: "list", input: "li42e2:hie", expected: test_list_simple, expectErr: false},
		{name: "list without close symbol", input: "li42e2:hi", expected: test_list_simple, expectErr: true},
		{name: "dictionary", input: "d6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42ee", expected: test_map_simple, expectErr: false},
		{name: "dictionary without close symbol", input: "d6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42e", expected: test_map_simple, expectErr: true},
		{name: "list with all nested data types", input: "li42e2:hid6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42eeli42e2:hiee", expected: test_list_complex, expectErr: false},
		{name: "dictionary with all nested data types", input: "d6:answeri42e5:hello5:world4:listli42e2:hie4:dictd6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42eee", expected: test_map_complex, expectErr: false},
		{name: "sequence of elements without dict or list wrap", input: "3:hi!i42e", expected: nil, expectErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := decodebencode.DecodeBencode(string(tc.input))
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil, result: %v", output)
				}
				if output != nil {
					t.Errorf("Expected nil result, got %v", output)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(output, tc.expected) {
					t.Errorf("Expected %v, got %v", tc.expected, output)
				}
			}
		})
	}

}
