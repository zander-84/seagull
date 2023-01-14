package tool

import (
	"fmt"
	"testing"
	"time"
)

func TestSliceArrayTime(t *testing.T) {
	//tn := time.Now()
	for i := 0; i < 100; i++ {
		ts, _ := defaultTime.SliceArrayTime(time.Now().AddDate(-1, 0, 0), time.Now(), time.Hour)

		fmt.Println(len(ts))
		//tag:=make(map[int64]bool,0)
		for _, v := range ts {
			if v[0].Unix()*1000 == v[1].Unix()*1000 {
				fmt.Println("err:", v[0].UnixNano(), v[1].UnixNano())
			}
			//fmt.Println(v[0].Format("2006-01-02 15:04:05"), v[1].Format("2006-01-02 15:04:05"))
		}

		fmt.Println()
	}
}

func TestSliceArraySecondTime(t *testing.T) {
	//tn := time.Now()
	//for i := 0; i < 100; i++ {
	//	ts, _ := defaultTime.SliceArraySecondTime(time.Now().AddDate(-1, 0, 0), time.Now(), time.Hour)
	//
	//	fmt.Println(len(ts))
	//	//tag:=make(map[int64]bool,0)
	//	for _, v := range ts {
	//		if v[0].Unix()*1000 == v[1].Unix()*1000 {
	//			fmt.Println("err:", v[0].UnixNano(), v[1].UnixNano())
	//		}
	//		//fmt.Println(v[0].Format("2006-01-02 15:04:05"), v[1].Format("2006-01-02 15:04:05"))
	//	}
	//
	//	fmt.Println()
	//}

	ts, _ := defaultTime.SliceArraySecondTime(time.Now().AddDate(0, 0, -3), time.Now(), time.Hour)

	fmt.Println(len(ts))
	//tag:=make(map[int64]bool,0)
	for _, v := range ts {
		if v[0].Unix()*1000 == v[1].Unix()*1000 {
			fmt.Println("err:", v[0].UnixNano(), v[1].UnixNano())
		}
		fmt.Println(v[0].Format("2006-01-02 15:04:05"), v[1].Format("2006-01-02 15:04:05"))
	}

	fmt.Println()
}
