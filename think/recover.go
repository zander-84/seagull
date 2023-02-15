package think

import (
	"context"
	"log"
	"runtime"
)

func Recover(ctx context.Context) error {
	if rErr := recover(); rErr != nil {
		buf := make([]byte, 64<<10)
		n := runtime.Stack(buf, false)
		buf = buf[:n]
		log.Printf("Printf err: %v \n", rErr)
		log.Println(string(buf))

		return nil
	}
	return nil
}
