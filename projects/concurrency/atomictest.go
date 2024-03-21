package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var x int32 = 0

func increment(wg *sync.WaitGroup) { //
	swapSuccess := false
	for !swapSuccess {
		curVal := x
		nextVal := curVal + 1
		swapSuccess = atomic.CompareAndSwapInt32(&x, curVal, nextVal)
	}
	wg.Done()
}

func main() {
	var w sync.WaitGroup
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go increment(&w)
	}
	w.Wait()
	fmt.Println("final value of x", x)
}
