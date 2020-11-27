package main

import (
	"context"
	"sync"
)

// ExecDel is the delegate for creating a task
type ExecDel func(context.Context) (interface{}, error)

// Status represent task execute status
type Status int

var (
	NotStarted Status = 0
	Running    Status
	Finish     Status
)

type TaskResult struct {
	Result interface{}
	Error  error
}

type Task struct {
	ctx     context.Context
	once    sync.Once
	execDel ExecDel
	wg      *sync.WaitGroup

	taskResult *TaskResult
	Status     Status
}

// Start will trigger task execute
func (t *Task) Start() {
	t.once.Do(t.coreExec)
}

// WaitForResult will return the result only after execution finish
func (t *Task) WaitForResult() *TaskResult {
	t.wg.Wait()
	return t.taskResult
}

func (t *Task) coreExec() {
	t.Status = Running
	go func(task *Task) {
		ret, err := t.execDel(t.ctx)
		t.taskResult = &TaskResult{
			Result: ret,
			Error:  err,
		}
		t.Status = Finish
		t.wg.Done()
	}(t)
}

// New used to create new task
func New(ctx context.Context, execDel ExecDel) *Task {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	return &Task{
		ctx:    ctx,
		Status: NotStarted,
		once:   sync.Once{},
		wg:     wg,
	}
}

// Run will create new task and start
func Run(ctx context.Context, execDel ExecDel) *Task {
	t := New(ctx, execDel)
	t.Start()
	return t
}
