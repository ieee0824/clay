package clay

import (
	"testing"
	"reflect"
)

func TestConvertMap(t *testing.T) {
	tests := []struct {
		input map[string]interface{}
		want map[string]interface{}
		err bool
	}{
		{
			input:map[string]interface{}{
				"slice": []string{"foo", "bar", "baz"},
				"map": map[string]interface{}{
					"country": "japan",
					"llustrator": []string{
						"カントク",
						"鈴平ひろ",
					},
				},
			},
			want: map[string]interface{}{
				"/slice": []string{"foo", "bar", "baz"},
				"/map/country": "japan",
				"/map/llustrator": []string{"カントク","鈴平ひろ"},
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := convertMap(test.input, "")
		if !test.err && err != nil {
			t.Fatalf("should not be error for %v but: %v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %v but not:", test.input)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}
