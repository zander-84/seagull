package contract

// CheckCode 校验码
type CheckCode interface {
	Check(in string, code string) error
	Sign(in string) string
}
