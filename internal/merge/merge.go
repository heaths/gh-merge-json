package merge

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/cli/go-gh/pkg/jsonpretty"
	"github.com/imdario/mergo"
)

func MergeJSON(r io.Reader, w io.Writer, tty bool) error {
	indent := ""
	if tty {
		indent = "  "
	}

	var merged interface{}

	dec := json.NewDecoder(r)
	for dec.More() {
		var data json.RawMessage

		// Consume the entire token.
		err := dec.Decode(&data)
		if err != nil {
			return err
		}

		buf, err := data.MarshalJSON()
		if err != nil {
			return err
		}

		err = merge(&merged, buf)
		if err != nil {
			return err
		}
	}

	buf, err := json.Marshal(merged)
	if err != nil {
		return err
	}

	r = bytes.NewBuffer(buf)
	return jsonpretty.Format(w, r, indent, tty)
}

func merge(dst *interface{}, buf []byte) error {
	var v interface{}
	if len(buf) == 0 {
		return errors.New("buffer is empty")
	} else if buf[0] == '{' {
		if *dst == nil {
			*dst = new(map[string]interface{})
		}
		v = new(map[string]interface{})
	} else if buf[0] == '[' {
		if *dst == nil {
			*dst = new([]interface{})
		}
		v = new([]interface{})
	}

	err := json.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	return mergo.Merge(*dst, v, mergo.WithAppendSlice, mergo.WithOverride)
}
