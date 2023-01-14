package codec

import (
	"encoding/xml"
)

var defaultXml = _xml{}

type _xml struct{}

// Marshal returns the wire format of v.
func (_xml) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// Unmarshal parses the wire format into v.
func (_xml) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// Name returns the name of the Codec implementation. The returned string
// will be used as part of content type in transmission.  The result must be
// static; the result cannot change between calls.
func (_xml) Name() string {
	return Xml
}
