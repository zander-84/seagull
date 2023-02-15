package codec

import (
	"github.com/zander-84/seagull/contract"
)

const Xml string = "xml"
const Json string = "json"
const Proto string = "proto"
const Yaml string = "Yaml"

var registeredCodecs = map[string]contract.Codec{
	Json:  defaultJson,
	Proto: defaultProto,
	Xml:   defaultXml,
	Yaml:  defaultYaml,
}

// RegisterCodec registers the provided Codec for use with all gRPC clients and
// servers.
//
// The Codec will be stored and looked up by result of its Name() method, which
// should match the content-subtype of the encoding handled by the Codec.  This
// is case-insensitive, and is stored and looked up as lowercase.  If the
// result of calling Name() is an empty string, RegisterCodec will panic. See
// Content-Type on
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests for
// more details.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe.  If multiple Codecs are
// registered with the same name, the one registered last will take effect.
func RegisterCodec(codec contract.Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	registeredCodecs[codec.Name()] = codec
}

// GetCodec gets a registered Codec by content-subtype, or nil if no Codec is
// registered for the content-subtype.
//
// The content-subtype is expected to be lowercase.
func GetCodec(contentSubtype string) contract.Codec {
	return registeredCodecs[contentSubtype]
}
