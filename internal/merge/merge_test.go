package merge

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{
			name: "object",
			in: `{
				"name": "foo",
				"values": [
					{
						"a": 1,
						"b": "foo"
					},
					{
						"a": 2,
						"b": "bar"
					}
				]
			}
			{
				"name": "bar",
				"values": [
					{
						"a": 2,
						"b": "baz"
					},
					{
						"a": 3,
						"b": "qux"
					}
				]
			}`,
			want: `{"name":"bar","values":[{"a":1,"b":"foo"},{"a":2,"b":"bar"},{"a":2,"b":"baz"},{"a":3,"b":"qux"}]}`,
		},
		{
			name: "array",
			in: `[
					{
						"a": 1,
						"b": "foo"
					},
					{
						"a": 2,
						"b": "bar"
					}
				]
				[
					{
						"a": 2,
						"b": "baz"
					},
					{
						"a": 3,
						"b": "qux"
					}
				]`,
			want: `[{"a":1,"b":"foo"},{"a":2,"b":"bar"},{"a":2,"b":"baz"},{"a":3,"b":"qux"}]`,
		},
		{
			name:    "invalid JSON",
			in:      `invalid JSON`,
			wantErr: true,
		},
		{
			name:    "invalid token",
			in:      `1234 invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.in)
			b, err := MergeJSON(r)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			buf := bytes.NewBuffer(b)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
