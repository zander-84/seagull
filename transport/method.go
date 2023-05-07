package transport

type Method string

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
)

func (m Method) String() string {
	return string(m)
}
