package util

import (
	"sync"
)

//如果结构体S，包含一个匿名字段T，那么这个结构体S 就有了T的方法。
type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1) //sync.WaitGroup结构中方法
	go func() {
		cb()
		w.Done() //sync.WaitGroup结构中方法
	}()
}
