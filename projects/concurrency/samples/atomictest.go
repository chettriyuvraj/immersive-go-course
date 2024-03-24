package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var atomicX atomic.Int32

func increment(wg *sync.WaitGroup) { //
	atomicX.Add(1)
	wg.Done()
}

func main() {
	var w sync.WaitGroup
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go increment(&w)
	}
	w.Wait()
	fmt.Println("final value of x", atomicX.Load())
}

/* Why did I create this convoluted implementation? */
// func increment(wg *sync.WaitGroup) { //
// 	for {
// 		curVal := atomicX.Load()
// 		if atomicX.CompareAndSwap(curVal, curVal+1) {
// 			break
// 		}
// 	}
// 	wg.Done()
// }

/* Alternate solution

var x int32 = 0

func increment(wg *sync.WaitGroup) {
	swapSuccess := false
	for !swapSuccess {
		curVal := x
		nextVal := curVal + 1
		swapSuccess = atomic.CompareAndSwapInt32(&x, curVal, nextVal)
	}
	wg.Done()
}

*/
