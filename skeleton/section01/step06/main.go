// STEP06: レア度ごとに出る確率を変えてみよう

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// 乱数の種を設定する
	// 現在時刻をUNIX時間にしたものを種とする
	rand.Seed(time.Now().Unix())

	// 0から99までの間で乱数を生成する
	num := rand.Intn(100)
	fmt.Println(num)
	// TODO: 変数numが0〜79のときは"ノーマル"、
	// 80〜94のときは"R"、95〜98のときは"SR"、
	// それ以外のときは"XR"と表示する
	switch {
	case num < 50:
		fmt.Println("N")
	case num < 70:
		fmt.Println("R")
	case num < 80:
		fmt.Println("SR")
	default:
		fmt.Println("SSR")
	}
}
