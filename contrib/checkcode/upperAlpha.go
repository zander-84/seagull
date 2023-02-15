package checkcode

import (
	"errors"
	"github.com/zander-84/seagull/contract"
	"hash/crc32"
)

type upperAlpha struct {
	size int
	code []string
}

func (t *upperAlpha) Check(in string, code string) error {
	if code == "" || len(code) != t.size {
		return errors.New("err data")
	}

	if code != t.Sign(in) {
		return errors.New("err data")
	}
	return nil
}

func (t *upperAlpha) Sign(in string) string {
	out := ""
	for i := t.size; i > 0; i-- {
		data := crc32.ChecksumIEEE([]byte(in))
		out += t.code[data%26]
		in += out
	}

	return out
}

func NewAlpha(size int) contract.CheckCode {
	out := new(upperAlpha)
	out.size = size
	out.code = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	return out
}
