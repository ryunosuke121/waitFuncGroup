# waitFuncGroup
ゴルーチンで実行したい関数をグループに追加すると、グループ内の関数がすべて実行完了するまで待機するライブラリです。
waitGroupは汎用性の高い仕様な一方で、初学者には難しいといった課題があります。
waitFuncGroupでは直感的な操作でゴルーチンを実行することができます。

# Getting Started/ スタートガイド
## パッケージのインストール
`go get github.com/ryunosuke121/waitFuncGroup`
## グループの作成
`wfg := waitFuncGroup.NewWaitFuncGroup(false)`
## タスクの追加
`wfg.Add([実行したい関数], "[タグ]")`
## タスクの終了を待機
`wfg.Wait()`

## タスクの監視
Wait()中に実行する関数をターミナルで監視することができます。どのタスクに時間がかかっているかデバッグに役立てることができます。

`wfg:= waitFuncGroup(true)`
または
`wfg.Monitor()`
### 監視モニター
https://github.com/ryunosuke121/waitFuncGroup/assets/117281628/8ebf38f1-81b3-4ba6-83fc-864ee5244abf

```
func main() {
	wfg := waitFuncGroup.NewWaitFuncGroup(false)
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
```
