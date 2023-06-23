package main

import (
	"time"

	"github.com/ryunosuke121/waitTask/waitFuncGroup"
)

func main() {
	wfg := waitFuncGroup.NewWaitFuncGroup()
	wfg.Monitor()
	wfg.Add(func() {
		time.Sleep(1 * time.Second)
	})
	wfg.Add(func() {
		time.Sleep(2 * time.Second)
		panic("panic")
	})
	wfg.Add(func() {
		time.Sleep(5 * time.Second)
	})
	wfg.Wait()
}
