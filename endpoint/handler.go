package endpoint

type Method string
type Kind string

const (
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodPatch   Method = "PATCH" // RFC 5789
	MethodDelete  Method = "DELETE"
	MethodOptions Method = "OPTIONS"

	MethodConnect Method = "CONNECT"
	MethodTrace   Method = "TRACE"

	Http   Kind = "HTTP"
	Grpc   Kind = "GRPC"
	Custom Kind = "CUSTOM"
	Empty  Kind = "EMPTY"
)

func (m Method) String() string {
	return string(m)
}

func (p Kind) IsHttp() bool {
	return p == Http
}

func (p Kind) IsEmpty() bool {
	return p == Empty
}

func (p Kind) IsGrpc() bool {
	return p == Grpc
}

func (p Kind) IsCustom() bool {
	return p == Custom
}

func inProtocols(s Kind, codec Codecs) bool {
	if codec == nil {
		return false
	}
	_, ok := codec[s]
	return ok
}
