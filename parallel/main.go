package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type runFunc func(context.Context) (interface{}, error)

func main() {
	fmt.Println("-----start 1----")
	baseCtx := context.Background()
	testRunAll(baseCtx)
	fmt.Println("------start 2----")
	ctx, cancel := context.WithTimeout(baseCtx, time.Millisecond*1200)
	defer cancel()
	testRunAll(ctx)
	fmt.Println("------start 3----")
	testRunWithTimeOut(baseCtx, time.Minute)
	fmt.Println("------start 4----")
	testRunWithTimeOut(baseCtx, time.Millisecond*500)
}

func testRunAll(ctx context.Context) {
	funcs := make([]runFunc, 10)
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(funcs); i++ {
		funcs[i] = func(ctx context.Context) (interface{}, error) {
			num := rand.Int() % 2000
			return run(ctx, num)
		}
	}

	start := time.Now().UnixNano()
	rets, errs := runParallel(ctx, funcs)
	parallelTime := time.Now().UnixNano() - start
	var maxItem int
	for i := range funcs {
		if item, ok := rets[i].(int); ok && item > maxItem {
			maxItem = item
		}

		fmt.Printf("ret %v, error %v\n", rets[i], errs[i])
	}

	fmt.Printf("cost time(ns): %v, maxItem: %v\n", parallelTime, maxItem)
}

func testRunWithTimeOut(ctx context.Context, timeout time.Duration) {
	funcs := make([]runFunc, 10)
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(funcs); i++ {
		funcs[i] = func(ctx context.Context) (interface{}, error) {
			num := rand.Int() % 2000
			return run(ctx, num)
		}
	}

	start := time.Now().UnixNano()
	rets, errs := runWithTimeout(ctx, funcs, timeout)
	parallelTime := time.Now().UnixNano() - start
	var maxItem int
	for i := range funcs {
		if item, ok := rets[i].(int); ok && item > maxItem {
			maxItem = item
		}

		fmt.Printf("ret %v, error %v\n", rets[i], errs[i])
	}

	fmt.Printf("cost time(ns): %v, maxItem: %v\n", parallelTime, maxItem)
}

func runWithTimeout(ctx context.Context, funcs []runFunc, timeout time.Duration) ([]interface{}, []error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	wg := sync.WaitGroup{}
	rets := make([]interface{}, len(funcs))
	errs := make([]error, len(funcs))
	wg.Add(len(funcs))
	for i, exeFunc := range funcs {
		go coreExec(ctx, i, exeFunc, &wg, rets, errs)
	}

	wg.Wait()
	return rets, errs
}

type retPack struct {
	Ret interface{}
	Err error
}

func coreExec(ctx context.Context, idx int, execFunc runFunc, wg *sync.WaitGroup, rets []interface{}, errs []error) {
	retChan := make(chan retPack, 1)
	go func() {
		ret, err := execFunc(ctx)
		retChan <- retPack{
			Ret: ret,
			Err: err,
		}
	}()

	select {
	case <-ctx.Done():
		errs[idx] = ctx.Err()
	case item := <-retChan:
		rets[idx], errs[idx] = item.Ret, item.Err
	}

	wg.Done()
}

func runParallel(ctx context.Context, funcs []runFunc) ([]interface{}, []error) {
	wg := sync.WaitGroup{}
	rets := make([]interface{}, len(funcs))
	errs := make([]error, len(funcs))
	wg.Add(len(funcs))
	for i, exeFunc := range funcs {
		go coreExec(ctx, i, exeFunc, &wg, rets, errs)
	}

	wg.Wait()
	return rets, errs
}

func run(ctx context.Context, num int) (int, error) {
	if num%10 == 1 {
		return 0, fmt.Errorf("invalid num %v", num)
	}

	c := make(chan struct{}, 1)
	go func() {
		time.Sleep(time.Duration(num) * time.Millisecond)
		c <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			return 0, err
		}
	case <-c:
	}

	return num, nil
}
