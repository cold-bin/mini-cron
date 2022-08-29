// @author cold bin
// @date 2022/7/19

package main

import (
	"time"
)

type Struct struct {
	A int
	B string
	C []byte
	D time.Time
}

func some(b uint8, a int, c string, d []byte, e Struct) (int, uint8, int, string, Struct) {
	e.A = a
	e.B = c
	e.D = time.Now()

	return 2333, uint8(11), 3222, "hello world", e
}

func main() {
	wheel := NewTimeWheel(12, 1, time.Second)
	go wheel.Start()
	////获取函数名和地址
	params := make([]interface{}, 0, 5)
	params = append(params, uint8(2))
	params = append(params, 1)
	params = append(params, "hello world")
	params = append(params, []byte("abcdef"))
	params = append(params, *new(Struct))

	wt := NewWorkTicker(1, some, params)
	wheel.AddWorkTicker(wt)
	//另起协程轮询时间轮
	time.Sleep(120 * time.Second)
}
