package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	multiListen()
}

func multiListen() {
	w := &sync.WaitGroup{}
	go func() {
		w.Add(1)
		time.Sleep(time.Second * 3)
		w.Done()
	}()
	fmt.Println("start")
	w.Wait()
	fmt.Println("2" + strconv.Itoa(0))
	w.Wait()
	fmt.Println("3" + strconv.Itoa(0))

	ch := make(chan struct{}, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		<-ch
		fmt.Println("from 1")
		wg.Done()
	}(&wg)

	go func(wg *sync.WaitGroup) {
		<-ch
		fmt.Println("from 2")
		wg.Done()
	}(&wg)
	ch <- struct{}{}
	wg.Wait()
}
