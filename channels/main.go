package main

import (
	"fmt"
	"time"
)

func worker(workerId int, data chan int) {
	for x := range data {
		fmt.Printf("worker %d received %d\n", workerId, x)
		time.Sleep(time.Second)
	}
}

func main() {
	ch := make(chan int)
	qtdWorkers := 100

	for i := 0; i < qtdWorkers; i++ {
		go worker(i, ch)
	}

	for i := 0; i < 1000; i++ {
		ch <- i
	}
}
