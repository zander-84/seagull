package unique

import (
	"fmt"
	"github.com/zander-84/seagull/contrib/checkcode"
	"testing"
)

func TestNew(t *testing.T) {
	u := New("LF", "", "aa", checkcode.NewAlpha(3))
	fmt.Println(u.ID())
	fmt.Println(u.Check("LFYPG458000001221230170433"))

}

func BenchmarkAssert(b *testing.B) {
	u := New("LF", "", "aa", checkcode.NewAlpha(2))

	for i := 0; i < b.N; i++ {
		//crc32.ChecksumIEEE([]byte("aaaaaaaaaaaaaaaaaaaaa"))
		u.ID()
	}
}
