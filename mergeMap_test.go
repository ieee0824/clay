package clay

import (
	"reflect"
	"testing"
)

func TestMergeMap(t *testing.T) {
	tests := []struct {
		src  map[string]interface{}
		dst  map[string]interface{}
		want map[string]interface{}
	}{
		{
			map[string]interface{}{
				"hoge": "huga",
			},
			map[string]interface{}{},
			map[string]interface{}{
				"hoge": "huga",
			},
		},
		{
			map[string]interface{}{
				"hoge": "huga",
				"foo":  "bar",
			},
			map[string]interface{}{},
			map[string]interface{}{
				"hoge": "huga",
				"foo":  "bar",
			},
		},
		{
			map[string]interface{}{
				"hoge": "huga",
			},
			map[string]interface{}{
				"foo": "bar",
			},
			map[string]interface{}{
				"hoge": "huga",
				"foo":  "bar",
			},
		},
		{
			map[string]interface{}{
				"hoge": "huga",
			},
			map[string]interface{}{
				"hoge": "baz",
				"foo":  "bar",
			},
			map[string]interface{}{
				"hoge": "baz",
				"foo":  "bar",
			},
		},
	}

	for _, test := range tests {
		mergeMap(test.src, test.dst)

		if !reflect.DeepEqual(test.src, test.want) {
			t.Fatalf("want %q, but %q:", test.want, test.src)
		}
	}
}
