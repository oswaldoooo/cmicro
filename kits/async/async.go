package async

import (
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/oswaldoooo/cmicro/types"
)

const (
	FORCE_CLOSED    = 1
	GRACEFUL_CLOSED = 2
)

// async package
type TaskStack struct {
	stack      *types.Stack
	lock       sync.Mutex
	insidelock sync.Mutex
	cond       *sync.Cond
	closed     uint8
}
type Task struct {
	F    any
	Args []any
}

type Strace struct {
	Stack string
	Val   any
}

func NewTask() *TaskStack {
	var ts TaskStack
	ts.cond = sync.NewCond(&ts.insidelock)
	ts.stack = &types.Stack{}
	return &ts
}
func (ts *TaskStack) AddTask(f any, args ...any) {
	if ts.closed > 0 {
		panic("task stack is closed")
	}
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.stack.PushBack(Task{F: f, Args: args})
	ts.cond.Signal()
}
func (ts *TaskStack) Run() {
	for {
		ts.cond.L.Lock()
		for ts.stack.Size() == 0 && ts.closed == 0 {
			ts.cond.Wait()
		}
		ts.cond.L.Unlock()
		if ts.closed == FORCE_CLOSED || (ts.closed == GRACEFUL_CLOSED && ts.stack.Size() == 0) {
			return
		}
		ts.lock.Lock()
		task := ts.stack.PopHead().(Task)
		ts.lock.Unlock()
		run(task.F, task.Args...)
	}
}
func (ts *TaskStack) Start() <-chan Strace {
	st := make(chan Strace, 1)
	go start(ts.Run, st)
	return st
}
func (ts *TaskStack) Close() error {
	ts.closed = FORCE_CLOSED
	ts.cond.Signal()
	return nil
}
func (ts *TaskStack) GracefulClose() {
	ts.closed = GRACEFUL_CLOSED
	ts.cond.Signal()
}
func start(f func(), sc chan<- Strace) {
	defer func() {
		if r := recover(); r != nil {
			sc <- Strace{Stack: string(debug.Stack()), Val: r}
			start(f, sc)
		} else {
			close(sc)
			return
		}
	}()
	f()
}
func run(f any, args ...any) {
	ftp := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)

	if len(args) != ftp.NumIn() {
		panic("args required")
	}
	argslen := len(args)
	var arg []reflect.Value = make([]reflect.Value, 0, argslen)
	for i := 0; i < argslen; i++ {
		arg = append(arg, reflect.ValueOf(args[i]).Convert(ftp.In(i)))
	}
	fval.Call(arg)
}
