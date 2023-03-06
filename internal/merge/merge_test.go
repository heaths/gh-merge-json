package merge

import (
	"bytes"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tty     bool
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
			// BUGBUG: https://github.com/cli/go-gh/issues/109
			want: heredoc.Doc(`{
				"name": "bar",
				"values": [
				{
				"a": 1,
				"b": "foo"
				},
				{
				"a": 2,
				"b": "bar"
				},
				{
				"a": 2,
				"b": "baz"
				},
				{
				"a": 3,
				"b": "qux"
				}
				]
				}
			`),
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
			// BUGBUG: https://github.com/cli/go-gh/issues/109
			want: heredoc.Doc(`[
				{
				"a": 1,
				"b": "foo"
				},
				{
				"a": 2,
				"b": "bar"
				},
				{
				"a": 2,
				"b": "baz"
				},
				{
				"a": 3,
				"b": "qux"
				}
				]
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.in)
			w := &bytes.Buffer{}
			err := MergeJSON(r, w, tt.tty)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, w.String())
		})
	}
}
