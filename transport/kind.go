package transport

type Kind string

const (
	Http   Kind = "HTTP"
	Grpc   Kind = "GRPC"
	Custom Kind = "CUSTOM"
	Empty  Kind = "EMPTY"
)

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
