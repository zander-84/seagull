package codec

import (
	"encoding/json"
)

var defaultJson = _json{}

type _json struct{}

// Marshal returns the wire format of v.
func (_json) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the wire format into v.
func (_json) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Name returns the name of the Codec implementation. The returned string
// will be used as part of content type in transmission.  The result must be
// static; the result cannot change between calls.
func (_json) Name() string {
	return Json
}
