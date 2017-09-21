package clay

import (
	"testing"
	"reflect"
	"github.com/aws/aws-sdk-go/aws"
	"encoding/json"
)

type T struct {
	Foo string
	Bar *string
	Struct C
	StructPtr *C
}

func (t T)String() string {
	bin, _ := json.MarshalIndent(t, "", "    ")
	return string(bin)
}

type C struct {
	Int int
	IntPtr *int
}

func TestMold(t *testing.T) {
	tests := []struct{
		moldData string
		structData T
		want T
		err  bool
	}{
		{
			moldData: "{}",
			structData: T{},
			want: T{},
			err: false,
		},
		{
			moldData: `{"Foo": "hoge"}`,
			structData: T{},
			want: T{
				Foo: "hoge",
			},
			err: false,
		},
		{
			moldData: `{"Foo": "hoge", "Bar": "huga"}`,
			structData: T{},
			want: T{
				Foo: "hoge",
				Bar: aws.String("huga"),
			},
			err: false,
		},
		{
			moldData: `
			{
				"Struct":{
					"Int": "10"
				}
			}`,
			structData: T{},
			want: T{
				Struct: C{Int: 10},
			},
		},
		{
			moldData: `
			{
				"StructPtr":{
					"Int": "10"
				}
			}`,
			structData: T{},
			want: T{
				StructPtr: &C{Int: 10},
			},
		},
	}

	for _, test := range tests {
		mold := map[string]interface{}{}

		if err := json.Unmarshal([]byte(test.moldData), &mold); err != nil {
			panic(err)
		}
		err := Mold(mold, &test.structData)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %v but: %v", test.moldData, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %v but not:", test.moldData)
		}
		if !reflect.DeepEqual(test.structData, test.want) {
			t.Fatalf("want %q, but %q:", test.want, test.structData)
		}
	}
}