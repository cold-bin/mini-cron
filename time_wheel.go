// @author cold bin
// @date 2022/7/18

package main

import (
	"errors"
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
	tw.Slots = make([]*WorkTicker, slotNum, slotNum)
	tw.SlotsNum = slotNum
	tw.StopSignal = make(chan bool)
	tw.UnitTime = unitTime
	tw.TickerWheel = time.NewTicker(unitTime)
	//初始化槽位
	for i := 0; i < len(tw.Slots); i++ {
		tw.Slots = append(tw.Slots, &WorkTicker{})
	}

	tw.CurSlot = 0
	return
}

func (tw *TimeWheel) AddWorkTicker(wt *WorkTicker) {
	//检验
	if wt.StartTime < 0 || wt.ExpireDuration < 0 || wt.StartTime >= wt.ExpireDuration {
		log.Println("time is wrong.")
		return
	}
	//选择槽位
	tickets := wt.ExpireDuration - wt.StartTime
	//类似拉链法选择槽位：当前槽位往后偏移
	slotLoc := (tw.CurSlot + tickets/tw.DurationPerSlot) % tw.SlotsNum
	//在槽位处添加定时器
	for i := 0; i < tw.SlotsNum; i++ {
		//此处拉链
		if i == slotLoc {
			tmp := tw.Slots[i]
			for tmp.Next != nil {
				tmp = tmp.Next
			}
			tmp.Next = wt
			wt.Prev = tmp
			break
		}
	}
}

// DelALLWorkTicker 删除在当前after个单位时间之后的的所有定时任务
func (tw *TimeWheel) DelALLWorkTicker(after int) {
	if after <= 0 {
		log.Println("the after time is wrong.")
		return
	}

	slotLoc := (tw.CurSlot + after) % tw.SlotsNum
	//在槽位处删除所有定时器
	for i := 0; i < tw.SlotsNum; i++ {
		//此处拉链
		if i == slotLoc {
			if tw.Slots[i].Next != nil {
				prev := tw.Slots[i]
				next := tw.Slots[i].Next
				prev.Next = nil
				next.Prev = nil
			}
			break
		}
	}
}

// DelOneWorkTicker 删除在当前after个单位时间之后的的某个定时任务
func (tw *TimeWheel) DelOneWorkTicker(after int, fMap map[string]reflect.Value) (err error) {
	if after <= 0 {
		err = errors.New("the after time is wrong")
		return
	}

	slotLoc := (tw.CurSlot + after) % tw.SlotsNum
	//取出需要找的定时器任务
	key1 := ""
	var value1 reflect.Value
	for key1 = range fMap {
		ok := false
		var v reflect.Value
		if v, ok = fMap[key1]; !ok {
			log.Println("not exist the key: ", v)
			return
		}
		value1 = v
	}
	//在槽位处删除指定定时器
	for i := 0; i < tw.SlotsNum; i++ {
		//此处拉链
		if i == slotLoc {
			tmp := tw.Slots[i]
			for tmp != nil {
				if v, ok := tmp.FuncMap[key1]; ok {
					if v == value1 {
						//找到该定时任务
						if tmp.Next != nil {
							tmp.Next = tmp.Prev
							tmp.Prev = tmp
							return
						}
						tmp.Prev.Next = nil
						tmp.Prev = nil
						return
					}
				}
				tmp = tmp.Next
			}
			break
		}
	}
	return errors.New("not found the work")
}

func (tw *TimeWheel) Stop() {
	tw.StopSignal <- true
}

// Start 该方法将轮询时间，阻塞式的
func (tw *TimeWheel) Start() {
	defer func() {
		close(tw.StopSignal)
		for _, v := range tw.Slots {
			tmp := v
			for tmp != nil {
				for key := range tmp.FuncMap {
					if _, ok := tmp.FuncMap[key]; ok {
						delete(tmp.FuncMap, key)
					}
				}
				tmp = tmp.Next
			}
		}
	}()

	for {
		select {
		case <-tw.StopSignal: //结束信号
			return
		case <-tw.TickerWheel.C: //时钟滴答一次，进一个槽位
			tw.CurSlot++
			tw.CurSlot = tw.CurSlot % tw.SlotsNum
			//todo 将当前定时器里的任务提取出来，依次执行
			for i := 0; i < tw.SlotsNum; i++ {
				wt := tw.Slots[i]
				tmp := wt
				for tmp != nil {
					values, err := tmp.Execute()
					if err != nil {
						log.Println(err)
					}
					//打印结果至终端
					log.Println(values)
					tmp = tmp.Next
				}
			}
		default:
			return
		}
	}
}

func (tw *TimeWheel) Reset() {
	tw.CurSlot = 0
	tw.StartTime = time.Now()
}
