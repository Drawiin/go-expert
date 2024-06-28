package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Message struct {
	Id      int64
	Payload string
}

func main() {
	ch := make(chan Message)
	ch2 := make(chan Message)
	var id int64

	go func() {
		for {
			time.Sleep(1 * time.Second)
			atomic.AddInt64(&id, 1)
			ch <- Message{Id: id, Payload: "Service 1"}
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			atomic.AddInt64(&id, 1)
			ch2 <- Message{Id: id, Payload: "Service 2"}
		}
	}()

	for {
		select {
		case data := <-ch:
			fmt.Printf("Received channel 1 = %v\n", data)
		case data := <-ch2:
			fmt.Printf("Received channel 2 =  %v\n", data)
		case <-time.After(2 * time.Second):
			println("timeout")
		}
	}
}
