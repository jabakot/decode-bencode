# decode-bencode

A Go lang library for encoding and decoding data in bencode format.

![Beencode](logo.png "bee encode")

## Get the pkg

```bash
go get -u github.com/jabakot/decode-bencode
```

## Import pkg

```go
import (
	decodebencode "github.com/jabakot/decode-bencode"
)
```

## Decode bencode

```go
input := "d6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42ee"

output, err := decodebencode.DecodeBencode(input)

if (err != nil) {
    fmt.Println(err)
}

fmt.Println(output)
// map[answer:42 hello:world hi:mark jojo:ゴゴゴゴ wrong-answer:-42]

```

## Encode int to bencode (str)

```go
bencode := decodebencode.EncodeBencodeInteger(42)
// bencode == "i42e"
```

## Encode str to bencode (str)


```go
bencode1 := decodebencode.EncodeBencodeString("hi!")
// bencode1 == "3:hi!"
bencode2 := decodebencode.EncodeBencodeString("ゴゴゴゴ")
// bencode2 == "12:ゴゴゴゴ"
```


## Encode slice/array to bencode (str)


```go
listWithIntAndString := make([]any, 2)
listWithIntAndString[0] = 1
listWithIntAndString[1] = "hi!"

bencode := decodebencode.EncodeBencodeList(listWithIntAndString)
// bencode == "li1e3:hi!e"
```

## Encode map to bencode (str)

```go
dictWithStringsAndInt := make(map[string]any)
dictWithStringsAndInt["answer"] = 42
dictWithStringsAndInt["hello"] = "world"
dictWithStringsAndInt["hi"] = "mark"
dictWithStringsAndInt["jojo"] = "ゴゴゴゴ"
dictWithStringsAndInt["wrong-answer"] = -42

bencode := decodebencode.EncodeBencodeDict(listWithIntAndString)
// bencode == "d6:answeri42e5:hello5:world2:hi4:mark4:jojo12:ゴゴゴゴ12:wrong-answeri-42ee"
```


