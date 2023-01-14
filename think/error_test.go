package think

import (
	"errors"
	"fmt"
	"testing"
)

func TestChain(t *testing.T) {

	e0 := errors.New("e0")

	e11 := ErrAlert("alter").WithCause(e0)

	e13 := ErrRecordNotFound("miss").WithCause(e11)

	fmt.Println(e13.Error())                               // miss
	fmt.Println(errors.Unwrap(e13).Error())                // alter
	fmt.Println(errors.Unwrap(errors.Unwrap(e13)).Error()) // e0
}
