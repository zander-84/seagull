package think

import "errors"

type Mode string

const (
	Prod  Mode = "prod"  // 生产
	Dev   Mode = "dev"   // 测试
	Local Mode = "local" // 本地
)

func (m Mode) IsProd() bool {
	if m == Prod {
		return true
	}
	return false
}

func (m Mode) IsDev() bool {
	if m == Dev {
		return true
	}
	return false
}

func (m Mode) IsLocal() bool {
	if m == Local {
		return true
	}
	return false
}

func NewMode(m string) (Mode, error) {
	if m == string(Prod) {
		return Prod, nil
	} else if m == string(Dev) {
		return Dev, nil
	} else if m == string(Local) {
		return Local, nil
	} else {
		return "", errors.New("错误类型")
	}
}
