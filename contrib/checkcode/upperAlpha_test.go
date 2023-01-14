package checkcode

import (
	"fmt"
	"testing"
)

func TestNewAlpha(t *testing.T) {
	a := NewAlpha(3)

	a1 := a.Sign("abc")
	fmt.Println(a1)

	if err := a.Check("abc", "EJH"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("success")
	}
}
