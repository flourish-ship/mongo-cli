package httplib

import (
	"encoding/json"
	"io"
	"net/url"
	"reflect"
)

// API params
type Params map[string]interface{}

// Encode params to query string
func (params Params) EncodeUrlParam(writer io.Writer) {
	if params == nil || len(params) == 0 {
		return
	}

	written := false

	for k, v := range params {
		if written {
			io.WriteString(writer, "&")
		}

		io.WriteString(writer, url.QueryEscape(k))
		io.WriteString(writer, "=")

		if reflect.TypeOf(v).Kind() == reflect.String {
			io.WriteString(writer, url.QueryEscape(reflect.ValueOf(v).String()))
		} else {
			jsonStr, err := json.Marshal(v)

			if err != nil {
				return
			}

			io.WriteString(writer, url.QueryEscape(string(jsonStr)))
		}

		written = true
	}
}

func (params Params) EncodeRequestBody(writer io.Writer) {
	b, err := json.Marshal(params)
	if err != nil {
		return
	}

	writer.Write(b)
}
