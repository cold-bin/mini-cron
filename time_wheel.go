// @author cold bin
// @date 2022/7/18

package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

// TimeWheel 定时器挂载时间轮
type TimeWheel struct {
	//单位时间
	UnitTime time.Duration
	//每隔多少单位时间转动一个槽位
	DurationPerSlot int
	//时间轮的开始时间
	StartTime time.Time
	//当前时间指针
	CurTime time.Time
	//当前指向时间槽的编号
	CurSlot int
	//指定数量的槽
	Slots []*WorkTicker
	//槽的数量
	SlotsNum int
	//轮询停止信号
	StopSignal chan bool
	//定时器基础
	TickerWheel *time.Ticker
}

func NewTimeWheel(slotNum, dPerSlot int, unitTime time.Duration) (tw *TimeWheel) {
	tw = new(TimeWheel)
	tw.DurationPerSlot = dPerSlot
	tw.StartTime = time.Now()
	tw.CurTime = time.Now()
	tw.Slots = make([]*WorkTicker, slotNum, slotNum*2)
	tw.SlotsNum = slotNum
	tw.StopSignal = make(chan bool)
	tw.UnitTime = unitTime
	tw.TickerWheel = time.NewTicker(unitTime)
	//初始化槽位
	for i := 0; i < slotNum; i++ {
		tw.Slots[i] = &WorkTicker{}
	}

	tw.CurSlot = 0
	return
}

func (tw *TimeWheel) AddWorkTicker(wt *WorkTicker) {
	//检验
	if wt.ExpireDuration < 0 {
		log.Println("time is wrong.")
		return
	}
	//选择槽位 类似拉链法选择槽位：当前槽位往后偏移
	slotLoc := (tw.CurSlot + wt.ExpireDuration/tw.DurationPerSlot) % tw.SlotsNum

	//在槽位处添加定时器，此处拉链
	tmp := tw.Slots[slotLoc]
	for tmp.Next != nil {
		tmp = tmp.Next
	}
	tmp.Next = wt
	wt.Prev = tmp
}

// DelALLWorkTicker 删除在当前after个单位时间之后的的所有定时任务
func (tw *TimeWheel) DelALLWorkTicker(after int) {
	if after <= 0 {
		log.Println("the after time is wrong.")
		return
	}

	slotLoc := (tw.CurSlot + after) % tw.SlotsNum
	//在槽位处删除所有定时器
	if tw.Slots[slotLoc].Next != nil {
		prev := tw.Slots[slotLoc]
		next := tw.Slots[slotLoc].Next
		prev.Next = nil
		next.Prev = nil
	}
}

// DelOneWorkTicker 删除在当前after个单位时间之后的的某个定时任务，该空接口接收一个函数类型
func (tw *TimeWheel) DelOneWorkTicker(after int, f interface{}) (err error) {
	if after <= 0 {
		err = errors.New("the after time is wrong")
		return
	}
	typeOfF := reflect.TypeOf(f)
	//如果不是函数类型，不支持定时任务
	if typeOfF.Kind() != reflect.Func {
		err = errors.New("the interface don't accepted the Func type")
		return
	}

	slotLoc := (tw.CurSlot + after) % tw.SlotsNum

	//在槽位处删除指定定时器
	tmp := tw.Slots[slotLoc].Next
	for tmp != nil {
		if reflect.DeepEqual(tmp.Func, f) {
			if tmp.Next == nil {
				tmp.Prev.Next = nil
				tmp.Prev = nil
				return
			}
			tmp.Next = tmp.Prev
			tmp.Prev = tmp
			return
		}
		tmp = tmp.Next
	}
	return errors.New("not found the work func")
}

func (tw *TimeWheel) Stop() {
	tw.StopSignal <- true
}

// Start 该方法将轮询时间，阻塞式的
func (tw *TimeWheel) Start() {
	defer func() {
		close(tw.StopSignal)
		tw.TickerWheel.Stop()
	}()

	for {
		select {
		case <-tw.StopSignal: //结束信号
			fmt.Println("stop over...")
			return
		case <-tw.TickerWheel.C: //时钟滴答一次，进一个槽位
			// 开个协程异步执行，增加时间轮的精度
			go func(tw *TimeWheel) {
				tw.CurSlot++
				tw.CurSlot = tw.CurSlot % tw.SlotsNum
				//将当前定时器的挂载链表取出来，不包括头节点
				wt := tw.Slots[tw.CurSlot]
				tmp := wt.Next
				for tmp != nil {
					values, err := tmp.Execute()
					if err != nil {
						fmt.Println(err)
					}
					//打印结果至终端
					for _, v := range values {
						fmt.Println("values: ", v.Int())
					}
					//fmt.Println("values: ", values[0].Int())
					tmp = tmp.Next
				}
				fmt.Println("嘀嗒...")
			}(tw)
		default:
		}
	}
}

func (tw *TimeWheel) Reset() {
	tw.CurSlot = 0
	tw.StartTime = time.Now()
	tw.TickerWheel.Reset(tw.UnitTime)
}
