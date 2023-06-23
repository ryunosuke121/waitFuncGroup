package waitFuncGroup

import (
	"sync"
	"time"

	"github.com/rivo/tview"
)

// 関数が渡されると、その関数をゴルーチンで実行する
// タイムスケジュール機能
// どれかのタスクがパニックを起こしたときにどうするかの検討
// -> そのタスクをキャンセルして、他のタスクを実行する or そのまますべて終了させる

type WaitFuncGroup struct {
	wg       sync.WaitGroup
	funcs    map[int]func()
	progress map[int]TaskStatus
	mu       sync.Mutex
	ch       chan TaskStatus
	monitor  bool
}

type TaskStatus struct {
	taskid int
	done   bool
	status bool
}

func (wfg *WaitFuncGroup) Add(f func()) {
	wfg.wg.Add(1)

	wfg.mu.Lock()
	func_id := len(wfg.funcs) + 1
	wfg.funcs[func_id] = f
	wfg.progress[func_id] = TaskStatus{taskid: func_id, done: false, status: false}
	wfg.mu.Unlock()

	go func(id int) {
		defer func() {
			if err := recover(); err != nil {
				wfg.wg.Done()
				wfg.ch <- TaskStatus{taskid: id, status: false, done: true}
				wfg.mu.Lock()
				for i := 1; i < len(wfg.progress)+1; i++ {
					if !wfg.progress[i].done {
						wfg.mu.Unlock()
						return
					}
				}
				close(wfg.ch)
				wfg.mu.Unlock()
			} else {
				wfg.wg.Done()
			}
		}()
		f()
		wfg.ch <- TaskStatus{taskid: id, status: true, done: true}
		wfg.mu.Lock()
		for i := 1; i < len(wfg.progress)+1; i++ {
			if !wfg.progress[i].done {
				wfg.mu.Unlock()
				return
			}
		}
		close(wfg.ch)
		wfg.mu.Unlock()
	}(func_id)
}

// タスクが完了するまで待つ
func (wfg *WaitFuncGroup) Wait() {
	var table *tview.Table
	app = tview.NewApplication()

	if wfg.monitor {
		table = createTable()
		go func() {
			wfg.mu.Lock()
			for i := 1; i < len(wfg.progress)+1; i++ {
				if wfg.progress[i].done {
					if wfg.progress[i].status {
						setCompleteRow(table, i)
					} else {
						setPanicRow(table, i)
					}
				} else {
					setWorkingRow(table, i)
				}
			}
			wfg.mu.Unlock()
		}()
	}
	go func() {
		for taskstatus := range wfg.ch {
			wfg.mu.Lock()
			wfg.progress[taskstatus.taskid] = taskstatus
			wfg.mu.Unlock()
			if wfg.monitor {
				if taskstatus.status {
					setCompleteRow(table, taskstatus.taskid)
				} else {
					setPanicRow(table, taskstatus.taskid)
				}
			}
		}
	}()
	if wfg.monitor {
		go func() {
			display(table)
		}()
	}
	wfg.wg.Wait()
	if wfg.monitor {
		time.Sleep(1 * time.Second)
		app.Stop()
	}
}

// 新しいWaitFuncGroupを作成する
func NewWaitFuncGroup() *WaitFuncGroup {
	return &WaitFuncGroup{
		funcs:    make(map[int]func()),
		progress: make(map[int]TaskStatus),
		ch:       make(chan TaskStatus),
		monitor:  false,
	}
}

// 実行中の関数を表示する
func (wfg *WaitFuncGroup) Monitor() {
	wfg.monitor = true
}
