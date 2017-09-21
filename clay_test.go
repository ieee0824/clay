package clay

import (
	"testing"
	"reflect"
	"github.com/aws/aws-sdk-go/aws"
)

func TestMold(t *testing.T) {
	type T struct {
		Foo string
		Bar *string
	}
	tests := []struct{
		moldData interface{}
		structData T
		want T
		err  bool
	}{
		{
			moldData: map[string]interface{}{},
			structData: T{},
			want: T{},
			err: false,
		},
		{
			moldData: map[string]interface{}{
				"Foo": "hoge",
			},
			structData: T{},
			want: T{
				Foo: "hoge",
			},
			err: false,
		},
		{
			moldData: map[string]interface{}{
				"Foo": "hoge",
				"Bar": "huga",
			},
			structData: T{},
			want: T{
				Foo: "hoge",
				Bar: aws.String("huga"),
			},
			err: false,
		},
	}

	for _, test := range tests {
		err := Mold(test.moldData, &test.structData)
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