package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ieee0824/clay"
)

type T struct {
	Str      string
	StrPtr   *string
	Int      int
	IntPtr   *int
	Int64    int64
	Int64Ptr *int64
	C        Children
	CPtr     *Children
	Array    []Children
}

type Children struct {
	Foo string
	Bar *int
}

func main() {
	m := map[string]interface{}{}
	JSON := `
	{
		"Str": "hoge",
		"StrPtr": "hoge ptr",
		"IntPtr": "-10",
		"CPtr":{
			"Foo": "baz"
		},
		"Array": [
			{"Foo": "hoge"}
		]
	}
	`

	t := T{}

	if err := json.Unmarshal([]byte(JSON), &m); err != nil {
		log.Fatalln(err)
	}

	if err := clay.Mold(m, &t); err != nil {
		log.Fatalln(err)
	}

	result, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(result))
}
