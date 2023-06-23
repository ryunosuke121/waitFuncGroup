package waitFuncTest

import (
	"testing"
	"time"

	"sync"

	"go.uber.org/goleak"

	"github.com/ryunosuke121/waitTask/waitFuncGroup"
)

func TestWaitFuncGroup(t *testing.T) {
	wfg := waitFuncGroup.NewWaitFuncGroup(false)
	wfg.Add(func() {
		time.Sleep(2 * time.Second)
		//fmt.Println("1が実行されました")
	}, "1秒待機追加")
	wfg.Add(func() {
		time.Sleep(2 * time.Second)
		//fmt.Println("2が実行されました")
	}, "2秒待機追加")
	wfg.Add(func() {
		time.Sleep(3 * time.Second)
		//fmt.Println("3が実行されました")
	}, "3秒待機追加")
	wfg.Wait()
}

func TestWaitFuncGroup2(t *testing.T) {
	wfg := waitFuncGroup.NewWaitFuncGroup(false)
	wfg.Add(func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(1 * time.Second)
			wg.Done()
		}()
		wg.Wait()
	}, "1秒待機")
	wfg.Add(func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(2 * time.Second)
			wg.Done()
		}()
		wg.Wait()
	}, "2秒待機")
	wfg.Add(func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(3 * time.Second)
			wg.Done()
		}()
		wg.Wait()
	}, "3秒待機")
	wfg.Wait()
}

func TestWaitFuncGroup3(t *testing.T) {
	wfg := waitFuncGroup.NewWaitFuncGroup(false)
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(1 * time.Second)
			}()
		}, "1秒待機")
	}, "1秒待機")
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(2 * time.Second)
			}()
		}, "2秒待機")
	}, "2秒待機")
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(3 * time.Second)
			}()
		}, "3秒待機")
	}, "3秒待機")
	wfg.Wait()
}

func TestLeak(t *testing.T) {
	wfg := waitFuncGroup.NewWaitFuncGroup(false)
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(3 * time.Second)
			}()
		}, "3秒待機")
	}, "3秒待機")
	wfg.Add(func() {
		wfg.Add(func() {
			go func() {
				time.Sleep(2 * time.Second)
			}()
		}, "2秒待機")
	}, "2秒待機")
	wfg.Wait()
	defer goleak.VerifyNone(t)
}
