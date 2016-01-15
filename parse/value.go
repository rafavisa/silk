package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type errValue []byte

func (e errValue) Error() string {
	return fmt.Sprintf("invalid value: %s (did you forget quotes?)", string(e))
}

type Value struct {
	Data interface{}
}

func (v Value) String() string {
	b, err := json.Marshal(v.Data)
	if err != nil {
		panic("silk: cannot marshal value: \"" + fmt.Sprintf("%v", v.Data) + "\": " + err.Error())
	}
	return string(b)
}

// Equal gets whether the Data and specified value are equal.
// Supports regexp values.
func (v Value) Equal(val interface{}) bool {
	// check to see if this is regex
	if str, ok := v.Data.(string); ok {
		if strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/") {
			// looks like regexp to me
			regex := regexp.MustCompile(str[1 : len(str)-1])
			// turn the value into a string
			valStr := fmt.Sprintf("%v", val)
			if regex.Match([]byte(valStr)) {
				return true
			}
		}
	}
	return v.Data == val
}

func ParseValue(src []byte) *Value {
	var v interface{}
	src = bytes.TrimSpace(src)
	if err := json.Unmarshal(src, &v); err != nil {
		return &Value{Data: string(src)}
	}
	return &Value{Data: v}
}
