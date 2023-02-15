package tool

import "math"

var base62 = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var flip = _flip(base62)
var defaultBase62 = NewBase62()

type Base62 struct {
}

func NewBase62() *Base62 {
	return new(Base62)
}

func GetBase62() *Base62 {
	return defaultBase62
}
func (b *Base62) Encode(num int) string {
	baseStr := ""
	for {
		if num <= 0 {
			break
		}
		i := num % 62
		baseStr += base62[i]
		num = (num - i) / 62
	}
	return baseStr
}

func (b *Base62) Decode(in string) int {
	rs := 0
	len1 := len(in)
	for i := 0; i < len1; i++ {
		rs += flip[string(in[i])] * int(math.Pow(62, float64(i)))
	}
	return rs
}

func _flip(s []string) map[string]int {
	f := make(map[string]int)
	for index, value := range s {
		f[value] = index
	}
	return f
}
