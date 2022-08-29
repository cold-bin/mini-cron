// @author cold bin
// @date 2022/8/29

package main

import (
	"sync"
	"time"
)

type CronPool struct {
	Pool *sync.Pool
}

type newFunc func(slotNum, dPerSlot int, unitTime time.Duration) (tw *TimeWheel)

// NewCronPool 创建一个 TimeWheel 的对象池，创建函数是干净的，没有放入任何复用的对象
func NewCronPool(newFunc newFunc, slotNum, dPerSlot int, unitTime time.Duration) (pool *CronPool) {
	return &CronPool{Pool: &sync.Pool{
		// 当对象池里无对象可用时，使用该函数获取
		New: func() interface{} {
			return newFunc(slotNum, dPerSlot, unitTime)
		},
	}}
}

func (p *CronPool) Get() *TimeWheel {
	return p.Pool.Get().(*TimeWheel)
}

func (p *CronPool) Put(tw *TimeWheel) {
	p.Pool.Put(tw)
}
