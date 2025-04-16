package decodebencode_test

import (
	"reflect"
	"testing"

	decodebencode "github.com/jabakot/decode-bencode"
)

func TestIsDictSymbol(t *testing.T) {

	type TestCase[I any] struct {
		input    I
		expected bool
	}

	testTable := []TestCase[any]{
		{input: decodebencode.DICT_MARKER, expected: true},
		{input: decodebencode.LIST_MARKER, expected: false},
		{input: 1, expected: false},
		{input: "1", expected: false},
		{input: nil, expected: false},
	}

	for _, test := range testTable {
		output := decodebencode.IsDictSymbol(test.input)

		if output != test.expected {
			t.Errorf("got %v, wanted %v\n", output, test.expected)
		}
	}
}

func TestIsListSymbol(t *testing.T) {
	type TestCase[I any] struct {
		input    I
		expected bool
	}
	testTable := []TestCase[any]{
		{input: decodebencode.LIST_MARKER, expected: true},
		{input: decodebencode.DICT_MARKER, expected: false},
		{input: 1, expected: false},
		{input: "1", expected: false},
		{input: nil, expected: false},
	}

	for _, test := range testTable {
		output := decodebencode.IsListSymbol(test.input)

		if output != test.expected {
			t.Errorf("got %v, wanted %v\n", output, test.expected)
		}
	}
}

func TestStackPush(t *testing.T) {
	stack := make(decodebencode.DataStack, 0, 1)

	stack.Push(1)

	if len(stack) != 1 && stack[0] != 1 {
		t.Errorf("failing to append 1 to stack")
	}

}

func TestEmptyStackPop(t *testing.T) {
	stack := make(decodebencode.DataStack, 0, 1)

	result, err := stack.Pop()

	if err.Error() != "data stack is empty" {
		t.Errorf("received unexpected error, expected `data stack is empty`, got: %v", err)
	}

	if result != nil {
		t.Errorf("received not nil for empty stack, expected nil, got: %v", result)
	}
}

func TestStackPop(t *testing.T) {

	stack := make(decodebencode.DataStack, 0, 1)

	stack.Push(1)

	result, err := stack.Pop()

	if result != 1 {
		t.Errorf("received unexpected result on Pop(), expected result = 1, got %v", result)
	}

	if err != nil {
		t.Errorf("received unexpected error, exected error = nil, got %v", err)
	}

	if len(stack) != 0 {
		t.Errorf("lenth of stacked not changed after Pop(), expected len(stack) == 0, got %d", len(stack))
	}

}

func TestShrinkDictionary(t *testing.T) {
	type TestCase struct {
		name        string
		input       []interface{}
		want        map[string]interface{}
		expectError bool
	}

	testTable := []TestCase{
		{
			name:  "Valid input with 2 key-value pairs",
			input: []interface{}{"value1", "key1", "value2", "key2"},
			want:  map[string]interface{}{"key1": "value1", "key2": "value2"},
		},
		{
			name:        "Odd number of elements",
			input:       []interface{}{"only", "two", "values"},
			expectError: true,
		},
		{
			name:        "Non-string key",
			input:       []interface{}{"val", 123},
			expectError: true,
		},
		{
			name:  "Empty input",
			input: []interface{}{},
			want:  map[string]interface{}{},
		},
		{
			name:  "Mixed types in values",
			input: []interface{}{42, "x", true, "y"},
			want:  map[string]interface{}{"x": 42, "y": true},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			output, err := decodebencode.ShrinkDictionary(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(output, tc.want) {
					t.Errorf("Expected map %v, got %v", tc.want, output)
				}
			}
		})
	}
}

func TestShrinkStack(t *testing.T) {

	type TestCase struct {
		name        string
		input       decodebencode.DataStack
		expected    decodebencode.DataStack
		expectError bool
	}

	testCases := []TestCase{
		{
			name: "Shrink to dict",
			input: decodebencode.DataStack{
				decodebencode.DICT_MARKER, "key1", "value1", "key2", "value2",
			},
			expected: decodebencode.DataStack{
				map[string]interface{}{"key1": "value1", "key2": "value2"},
			},
		},
		{
			name: "Shrink to list",
			input: decodebencode.DataStack{
				decodebencode.LIST_MARKER, "item1", "item2",
			},
			expected: decodebencode.DataStack{
				[]interface{}{"item1", "item2"},
			},
		},
		{
			name: "Non-symbol tail, leaves last item",
			input: decodebencode.DataStack{
				"onlyOne",
			},
			expected: decodebencode.DataStack{
				"onlyOne",
			},
		},
		{
			name: "Shrink to dict with non-string key",
			input: decodebencode.DataStack{
				decodebencode.DICT_MARKER, 123, "value",
			},
			expectError: true,
		},
		{
			name: "Shrink to dict with odd elements",
			input: decodebencode.DataStack{
				decodebencode.DICT_MARKER, "value1", "key1", "extra",
			},
			expectError: true,
		},
		{
			name:     "Empty stack",
			input:    decodebencode.DataStack{},
			expected: decodebencode.DataStack{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stack := tc.input
			err := decodebencode.ShrinkStack(&stack)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil %v", stack...)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(stack, tc.expected) {
					t.Errorf("Expected stack %v, got %v", tc.expected, stack)
				}
			}
		})
	}
}
