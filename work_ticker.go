// @author cold bin
// @date 2022/7/18

package main

import (
	errors "errors"
	"reflect"
)

// WorkTicker 定时器
type WorkTicker struct {
	//记录当前定时器位于时间轮的哪个槽
	TimeSlot int
	//定时器任务存储：存储函数地址和函数名，使用时，可以通过反射调用函数（前提是必须显式的在项目里存在该函数），存一个函数
	FuncMap map[string]reflect.Value
	//函数的参数：需要使用reflect.ValueOf包裹
	FunParams []reflect.Value
	//指向下一个定时器
	Next *WorkTicker
	//指向上一个定时器任务
	Prev *WorkTicker
	//定时任务需要执行的时候，单位时间个数
	ExpireDuration int
	//定时器任务开始的时间
	StartTime int
}

// NewWorkTicker 创建一个新的定时器任务，此时的定时器任务是零散的，还没有添加到时间轮
func NewWorkTicker(sTime int, eTime int, fMap map[string]reflect.Value, fParams []reflect.Value) (wt *WorkTicker) {
	wt = new(WorkTicker)
	wt.ExpireDuration = eTime
	wt.StartTime = sTime
	wt.FuncMap = make(map[string]reflect.Value, 1)
	wt.FunParams = make([]reflect.Value, 0)
	//浅拷贝
	for key := range fMap {
		if v, ok := fMap[key]; ok {
			wt.FuncMap[key] = v
		}
	}
	wt.FunParams = append(wt.FunParams, fParams...)
	return
}

// Execute 执行任务
func (wt *WorkTicker) Execute() ([]reflect.Value, error) {
	for key := range wt.FuncMap {
		if f, ok := wt.FuncMap[key]; ok {
			if len(wt.FunParams) != f.Type().NumIn() {
				return nil, errors.New("the number of input params not match")
			}
			in := make([]reflect.Value, len(wt.FunParams))
			for k, v := range wt.FunParams {
				in[k] = reflect.ValueOf(v)
			}
			return f.Call(in), nil
		}
	}
	return nil, errors.New("don't know the error")
}
