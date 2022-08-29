// @author cold bin
// @date 2022/7/18

package main

import (
	"errors"
	"log"
	"reflect"
)

// WorkTicker 定时器
type WorkTicker struct {
	//据当前时间，指定时间段后的过期时间
	ExpireDuration int
	//记录当前定时器位于时间轮的哪个槽
	TimeSlot int
	//函数的参数：需要使用reflect.ValueOf包裹
	FunParams []interface{}
	//定时器任务存储：存储函数地址和函数名，使用时，可以通过反射调用函数（前提是必须显式的在项目里存在该函数），存一个函数
	Func interface{}
	//指向下一个定时器
	Next *WorkTicker
	//指向上一个定时器任务
	Prev *WorkTicker
}

// NewWorkTicker 创建一个新的定时器任务，此时的定时器任务是零散的，还没有添加到时间轮
func NewWorkTicker(aDuration int, f interface{}, fParams []interface{}) (wt *WorkTicker) {
	wt = new(WorkTicker)
	wt.ExpireDuration = aDuration
	wt.FunParams = make([]interface{}, 0)
	wt.Func = f

	wt.FunParams = append(wt.FunParams, fParams...)
	return
}

// Execute 执行任务
func (wt *WorkTicker) Execute() ([]reflect.Value, error) {

	//取出函数地址
	f := reflect.ValueOf(wt.Func)
	if f.Type().Kind() != reflect.Func {
		return nil, errors.New("workTicker do not contain func type")
	}
	//取出函数的参数
	log.Println("参数个数：", len(wt.FunParams), f.Type().NumIn())
	if len(wt.FunParams) != f.Type().NumIn() {
		return nil, errors.New("the number of input params not match")
	}
	in := make([]reflect.Value, len(wt.FunParams))
	for k, v := range wt.FunParams {
		in[k] = reflect.ValueOf(v)
	}
	return f.Call(in), nil
}
