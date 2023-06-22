package waitFuncGroup

// import (
// 	"internal/race"
// 	"sync/atomic"
// 	"unsafe"
// )

// type WaitFuncGroup struct {
// 	state atomic.Uint64 // 上位32ビットが待ち合わせカウンター、下位32ビットが待機中のゴルーチン数
// 	sema  uint32
// }

// // 関数が渡されると、その関数をゴルーチンで実行する
// func (wfg *WaitFuncGroup) Add(f func()) {
// 	if race.Enabled {
// 		race.Disable()
// 		defer race.Enable()
// 	}
// 	state := wfg.state.Add(uint64(1) << 32)
// 	// 待ち合わせカウンター
// 	v := int32(state >> 32)
// 	// 待機中のゴルーチン数
// 	w := uint32(state)
// 	if race.Enabled && v == int32(1) {
// 		// The first increment must be synchronized with Wait.
// 		// Need to model this as a read, because there can be
// 		// several concurrent wg.counter transitions from 0.
// 		race.Read(unsafe.Pointer(&wfg.sema))
// 	}
// 	// カウンターが負の場合パニック
// 	if v < 0 {
// 		panic("sync: negative WaitGroup counter")
// 	}
// 	// ゴルーチンが待機中の場合
// 	if w != 0 && v == int32(1) {
// 		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
// 	}
// 	if v > 0 || w == 0 {
// 		return
// 	}

// 	if wfg.state.Load() != state {
// 		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
// 	}
// 	// Reset waiters count to 0.
// 	wfg.state.Store(0)
// 	for ; w != 0; w-- {
// 		runtime_Semrelease(&wg.sema, false, 0)
// 	}
// }

// // Done decrements the WaitGroup counter by one.
// func (wg *WaitFuncGroup) Done() {
// 	wg.Add(-1)
// }

// // Wait blocks until the WaitGroup counter is zero.
// func (wg *WaitGroup) Wait() {
// 	if race.Enabled {
// 		race.Disable()
// 	}
// 	for {
// 		state := wg.state.Load()
// 		v := int32(state >> 32)
// 		w := uint32(state)
// 		if v == 0 {
// 			// Counter is 0, no need to wait.
// 			if race.Enabled {
// 				race.Enable()
// 				race.Acquire(unsafe.Pointer(wg))
// 			}
// 			return
// 		}
// 		// Increment waiters count.
// 		if wg.state.CompareAndSwap(state, state+1) {
// 			if race.Enabled && w == 0 {
// 				// Wait must be synchronized with the first Add.
// 				// Need to model this is as a write to race with the read in Add.
// 				// As a consequence, can do the write only for the first waiter,
// 				// otherwise concurrent Waits will race with each other.
// 				race.Write(unsafe.Pointer(&wg.sema))
// 			}
// 			runtime_Semacquire(&wg.sema)
// 			if wg.state.Load() != 0 {
// 				panic("sync: WaitGroup is reused before previous Wait has returned")
// 			}
// 			if race.Enabled {
// 				race.Enable()
// 				race.Acquire(unsafe.Pointer(wg))
// 			}
// 			return
// 		}
// 	}
// }
