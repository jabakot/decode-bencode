package decodebencode_test

import (
	"testing"

	decodebencode "github.com/jabakot/decode-bencode"
)

type EncoderTestCase[I any] struct {
	input    I
	expected string
}

func tableRunner[I any](testTable []EncoderTestCase[I], encoder func(input I) string, t *testing.T) {
	for _, test := range testTable {

		output := encoder(test.input)

		if output != test.expected {
			t.Errorf("got %q, wanted %q\n", output, test.expected)
		}
	}
}

func TestEncodeInteger(t *testing.T) {

	testTable := []EncoderTestCase[int]{
		{input: 1, expected: "i1e"},
		{input: -2, expected: "i-2e"},
		{input: 0, expected: "i0e"},
		{input: -0, expected: "i0e"},
		{input: 000000, expected: "i0e"},
		{input: 0x00001, expected: "i1e"},
	}

	tableRunner(testTable, decodebencode.EncodeBencodeInteger, t)

}
func TestEncodeString(t *testing.T) {
	testTable := []EncoderTestCase[string]{
		{input: "", expected: ""},
		{input: "hi!", expected: "3:hi!"},
		{input: "ゴゴゴゴ", expected: "12:ゴゴゴゴ"},
		{input: "42", expected: "2:42"},
	}

	tableRunner(testTable, decodebencode.EncodeBencodeString, t)
}

func TestEncodeList(t *testing.T) {

	// li1e3:hi!e

	listWithIntAndString := make([]any, 2)
	// i1e
	listWithIntAndString[0] = 1
	// 3:hi!
	listWithIntAndString[1] = "hi!"

	listWithLists := make([]any, 2)
	listWithLists[0] = listWithIntAndString
	listWithLists[1] = listWithIntAndString

	// d5:hello5:world2:hi4:marke
	dict1 := make(map[string]any)
	dict1["hello"] = "world"
	dict1["hi"] = "mark"

	// d4:jojo12:ゴゴゴゴe
	dict2 := make(map[string]any)
	//    4:jojo	12:ゴゴゴゴ
	dict2["jojo"] = "ゴゴゴゴ"

	// ld5:hello5:world2:hi4:marked4:jojo12:ゴゴゴゴee
	listWithDicts := make([]any, 2)
	listWithDicts[0] = dict1
	listWithDicts[1] = dict2

	// li1e3:hi!li1e3:hi!eld5:hello5:world2:hi4:marked4:jojo12:ゴゴゴゴeee
	listWithMixed := make([]any, 4)
	listWithMixed[0] = 1
	listWithMixed[1] = "hi!"
	listWithMixed[2] = listWithIntAndString
	listWithMixed[3] = listWithDicts

	testTable := []EncoderTestCase[[]any]{
		{input: make([]any, 0), expected: "le"},
		{input: listWithIntAndString, expected: "li1e3:hi!e"},
		{input: listWithLists, expected: "lli1e3:hi!eli1e3:hi!ee"},
		{input: listWithDicts, expected: "ld5:hello5:world2:hi4:marked4:jojo12:ゴゴゴゴee"},
		{input: listWithMixed, expected: "li1e3:hi!li1e3:hi!eld5:hello5:world2:hi4:marked4:jojo12:ゴゴゴゴeee"},
	}

	tableRunner(testTable, decodebencode.EncodeBencodeList, t)
}

func TestEncodeDict(t *testing.T) {

	// d6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42ee
	dictWithStringsAndInt := make(map[string]any)
	dictWithStringsAndInt["answer"] = 42
	dictWithStringsAndInt["hello"] = "world"
	dictWithStringsAndInt["hi"] = "mark"
	dictWithStringsAndInt["jojo"] = "ゴゴゴゴ"
	dictWithStringsAndInt["wrong-answer"] = -42

	// li1e3:hi!e
	listWithIntAndString := make([]any, 2)
	// i1e
	listWithIntAndString[0] = 1
	// 3:hi!
	listWithIntAndString[1] = "hi!"

	// li0ei1ei2ei3ee
	listWithInts := make([]any, 4)
	for i := range len(listWithInts) {
		listWithInts[i] = i
	}

	// d4:listli1e3:hi!e7:numbersli0ei1ei2ei3eee
	dictWithLists := make(map[string]any)
	dictWithLists["list"] = listWithIntAndString
	dictWithLists["numbers"] = listWithInts

	// d4:dictd4:listli1e3:hi!e7:numbersli0ei1ei2ei3eeee
	dictWithDict := make(map[string]any)
	dictWithDict["dict"] = dictWithLists

	// d6:answeri42e4:dictd4:dictd4:listli1e3:hi!e7:numbersli0ei1ei2ei3eeee4:jojo12:ゴゴゴゴ4:listli1e3:hi!ee
	dictWithMixed := make(map[string]any)
	dictWithMixed["answer"] = 42
	dictWithMixed["dict"] = dictWithDict
	dictWithMixed["jojo"] = "ゴゴゴゴ"
	dictWithMixed["list"] = listWithIntAndString

	testTable := []EncoderTestCase[map[string]any]{
		{input: make(map[string]any), expected: "de"},
		{input: dictWithStringsAndInt, expected: "d6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42ee"},
		{input: dictWithLists, expected: "d4:listli1e3:hi!e7:numbersli0ei1ei2ei3eee"},
		{input: dictWithDict, expected: "d4:dictd4:listli1e3:hi!e7:numbersli0ei1ei2ei3eeee"},
		{input: dictWithMixed, expected: "d6:answeri42e4:dictd4:dictd4:listli1e3:hi!e7:numbersli0ei1ei2ei3eeee4:jojo12:ゴゴゴゴ4:listli1e3:hi!ee"},
	}

	tableRunner(testTable, decodebencode.EncodeBencodeDict, t)
}
