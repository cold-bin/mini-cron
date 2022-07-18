// @author cold bin
// @date 2022/7/19

package main

import (
	"fmt"
	"reflect"
	"time"
)

type ControllerMapsType map[string]reflect.Value

var ControllerMaps ControllerMapsType

type Routers struct {
}

func (this *Routers) Login(msg string) {
	fmt.Println("Login:", msg)
}

func (this *Routers) ChangeName(msg *string) {
	fmt.Println("ChangeName:", *msg)
	*msg = *msg + " Changed"
}

func main() {
	wheel := NewTimeWheel(60, 1, time.Second)
	//另起协程轮询
	go wheel.Start()

	var ruTest Routers
	crMap := make(ControllerMapsType, 0)
	vf := reflect.ValueOf(&ruTest)
	vft := vf.Type()
	//读取方法数量
	mNum := vf.NumMethod()
	fmt.Println("NumMethod:", mNum)

	//遍历路由器的方法，并将其存入控制器映射变量中
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		fmt.Println("index:", i, " MethodName:", mName)
		crMap[mName] = vf.Method(i)
	}
	params := []reflect.Value{reflect.ValueOf("test the handle")}

	NewWorkTicker(2, 4, crMap, params)
}
