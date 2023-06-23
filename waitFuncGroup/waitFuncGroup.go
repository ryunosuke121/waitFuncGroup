package waitFuncGroup

import (
	"sync"
	"time"

	"fmt"

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
	name   string
}

// タスクを追加する
func (wfg *WaitFuncGroup) Add(f func(), name string) {
	wfg.wg.Add(1)

	wfg.mu.Lock()
	func_id := len(wfg.funcs) + 1
	wfg.funcs[func_id] = f
	wfg.progress[func_id] = TaskStatus{taskid: func_id, done: false, status: false, name: name}
	wfg.mu.Unlock()

	go func(id int, name string) {
		defer func() {
			// panicが起きたときの処理
			if err := recover(); err != nil {
				fmt.Println(err)
				wfg.wg.Done()
				wfg.ch <- TaskStatus{taskid: id, status: false, done: true, name: name}
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
		// 正常終了時の処理
		wfg.ch <- TaskStatus{taskid: id, status: true, done: true, name: name}
		wfg.mu.Lock()
		for i := 1; i < len(wfg.progress)+1; i++ {
			if !wfg.progress[i].done {
				wfg.mu.Unlock()
				return
			}
		}
		close(wfg.ch)
		wfg.mu.Unlock()
	}(func_id, name)
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
						setCompleteRow(table, i, wfg.progress[i].name)
					} else {
						setPanicRow(table, i, wfg.progress[i].name)
					}
				} else {
					setWorkingRow(table, i, wfg.progress[i].name)
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
					setCompleteRow(table, taskstatus.taskid, taskstatus.name)
				} else {
					setPanicRow(table, taskstatus.taskid, taskstatus.name)
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
func NewWaitFuncGroup(monitor bool) *WaitFuncGroup {
	return &WaitFuncGroup{
		funcs:    make(map[int]func()),
		progress: make(map[int]TaskStatus),
		ch:       make(chan TaskStatus),
		monitor:  monitor,
	}
}

// 実行中の関数を表示する
func (wfg *WaitFuncGroup) Monitor() {
	wfg.monitor = true
}
