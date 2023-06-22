package waitFuncTest

import (
	"testing"
	"time"

	"sync"

	"github.com/ryunosuke121/waitTask/waitFuncGroup"
)

func TestWaitFuncGroup(t *testing.T) {
	wfg := waitFuncGroup.NewWaitFuncGroup()
	wfg.Add(func() {
		time.Sleep(2 * time.Second)
		//fmt.Println("1が実行されました")
	})
	wfg.Add(func() {
		time.Sleep(2 * time.Second)
		//fmt.Println("2が実行されました")
	})
	wfg.Add(func() {
		time.Sleep(3 * time.Second)
		//fmt.Println("3が実行されました")
	})
	wfg.Wait()
}

func TestWaitFuncGroup2(t *testing.T) {
	wfg := waitFuncGroup.NewWaitFuncGroup()
	wfg.Add(func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(1 * time.Second)
			wg.Done()
		}()
		wg.Wait()
	})
	wfg.Add(func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(2 * time.Second)
			wg.Done()
		}()
		wg.Wait()
	})
	wfg.Add(func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(3 * time.Second)
			wg.Done()
		}()
		wg.Wait()
	})
	wfg.Wait()
}

func TestWaitFuncGroup3(t *testing.T) {
	wfg := waitFuncGroup.NewWaitFuncGroup()
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(1 * time.Second)
			}()
		})
	})
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(2 * time.Second)
			}()
		})
	})
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(3 * time.Second)
			}()
		})
	})
	wfg.Wait()
}
