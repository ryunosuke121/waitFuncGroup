package waitFuncGroup

import (
	"fmt"
	"sync"
)

// 関数が渡されると、その関数をゴルーチンで実行する
// タイムスケジュール機能
// どれかのタスクがパニックを起こしたときにどうするかの検討

type WaitFuncGroup struct {
	wg       sync.WaitGroup
	funcs    map[int]func()
	progress map[int]bool
	mu       sync.Mutex
	ch       chan int
}

func (wfg *WaitFuncGroup) Add(f func()) {
	wfg.wg.Add(1)

	wfg.mu.Lock()
	func_id := len(wfg.funcs) + 1
	wfg.funcs[func_id] = f
	wfg.progress[func_id] = false

	go func(id int) {
		defer wfg.wg.Done()
		f()
		wfg.ch <- id
	}(func_id)
	wfg.mu.Unlock()
}

func (wfg *WaitFuncGroup) Done() {
	wfg.wg.Done()
}

func (wfg *WaitFuncGroup) Wait() {
	fmt.Println("実行中の関数")
	for k := range wfg.funcs {
		fmt.Println(k)
	}

	go func() {
		for id := range wfg.ch {
			wfg.mu.Lock()
			wfg.progress[id] = true
			wfg.mu.Unlock()
			fmt.Println("------------------------")
			for i := 1; i < len(wfg.progress)+1; i++ {
				if wfg.progress[i] {
					fmt.Printf("task%d:Completed\n", i)
				} else {
					fmt.Printf("task%d:Working\n", i)
				}
			}
			fmt.Println("------------------------")
		}
	}()
	wfg.wg.Wait()
	close(wfg.ch)
}

// 新しいWaitFuncGroupを作成する
func NewWaitFuncGroup() *WaitFuncGroup {
	return &WaitFuncGroup{
		funcs:    make(map[int]func()),
		progress: make(map[int]bool),
		ch:       make(chan int),
	}
}
