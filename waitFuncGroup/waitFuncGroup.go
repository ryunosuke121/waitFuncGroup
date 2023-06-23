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
	progress map[int]bool
	mu       sync.Mutex
	ch       chan int
	monitor  bool
}

func (wfg *WaitFuncGroup) Add(f func()) {
	wfg.wg.Add(1)

	wfg.mu.Lock()
	func_id := len(wfg.funcs) + 1
	wfg.funcs[func_id] = f
	wfg.progress[func_id] = false
	wfg.mu.Unlock()

	go func(id int) {
		defer wfg.wg.Done()
		f()
		wfg.ch <- id
		wfg.mu.Lock()
		for i := 1; i < len(wfg.progress)+1; i++ {
			if !wfg.progress[i] {
				wfg.mu.Unlock()
				return
			}
		}
		close(wfg.ch)
		wfg.mu.Unlock()
	}(func_id)
}

func (wfg *WaitFuncGroup) Done() {
	wfg.wg.Done()
}

func (wfg *WaitFuncGroup) Wait() {
	var table *tview.Table
	app = tview.NewApplication()

	if wfg.monitor {
		table = createTable()
		go func() {
			wfg.mu.Lock()
			for i := 1; i < len(wfg.progress)+1; i++ {
				if wfg.progress[i] {
					completeRow(table, i)
				} else {
					workingRow(table, i)
				}
			}
			wfg.mu.Unlock()
		}()
	}
	go func() {
		for id := range wfg.ch {
			wfg.mu.Lock()
			wfg.progress[id] = true
			wfg.mu.Unlock()
			if wfg.monitor {
				completeRow(table, id)
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
		progress: make(map[int]bool),
		ch:       make(chan int),
		monitor:  false,
	}
}

// 実行中の関数を表示する
func (wfg *WaitFuncGroup) Monitor() {
	wfg.monitor = true
}
