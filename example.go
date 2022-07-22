// @author cold bin
// @date 2022/7/19

package main

import (
	"fmt"
	"reflect"
	"time"
)

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func old(a string) string {
	return "a()=" + a
}

func main() {
	wheel := NewTimeWheel(12, 1, time.Second)
	fmt.Println("wheel: ", wheel)
	go wheel.Start()

	//获取函数名和地址
	funcV1 := reflect.ValueOf(add)
	//funName := runtime.FuncForPC(reflect.ValueOf(add).Pointer()).Name()
	//funName := funcV1.Type().Name()
	//fmt.Println("函数信息：", funcV1, funName)
	//获取函数参数
	in := funcV1.Type().NumIn()
	params := make([]interface{}, 0)
	for i := 0; i < in; i++ {
		params = append(params, i)
	}
	fmt.Println("params: ", params)
	wt := NewWorkTicker(1, add, params)
	fmt.Println("workerTicker: ", wt)

	wheel.AddWorkTicker(wt)
	fmt.Println("wheel: ", wheel)
	//另起协程轮询时间轮
	time.Sleep(120 * time.Second)
}
