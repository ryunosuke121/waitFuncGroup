package main

import (
	"time"

	"github.com/ryunosuke121/waitTask/waitFuncGroup"
)

func main() {
	wfg := waitFuncGroup.NewWaitFuncGroup(true)
	wfg.Monitor()
	wfg.Add(func() {
		time.Sleep(1 * time.Second)
	}, "1秒待機")
	wfg.Add(func() {
		time.Sleep(2 * time.Second)
		panic("panic!!")
	}, "2秒待機")
	wfg.Add(func() {
		time.Sleep(5 * time.Second)
	}, "5秒待機")
	wfg.Wait()
}
