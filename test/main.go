package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var tOperations int32

func inc() {
	atomic.AddInt32(&tOperations, 1)
	//tOperations++
}

func main() {

	for i := 0; i < 1000; i++ {
		go inc()
	}

	time.Sleep(2 * time.Millisecond)

	fmt.Println(tOperations)
}
