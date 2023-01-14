package codec

import (
	"gopkg.in/yaml.v3"
)

var defaultYaml = _yaml{}

type _yaml struct{}

// Marshal returns the wire format of v.
func (_yaml) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal parses the wire format into v.
func (_yaml) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// Name returns the name of the Codec implementation. The returned string
// will be used as part of content type in transmission.  The result must be
// static; the result cannot change between calls.
func (_yaml) Name() string {
	return Yaml
}
