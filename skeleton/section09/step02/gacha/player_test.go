// TODO: gachaパッケージとは別でgacha_testパッケージにする
package gacha_test

import (
	"testing"

	"gacha/skeleton/section09/step02/gacha"
)

func TestPlayer_DrawableNum(t *testing.T) {
	cases := map[string]struct {
		tickets int
		coin    int
		want    int
	}{
		"zero-zero": {0, 0, 0},
		"plus-zero": {10, 0, 10},
		"plus-plus": {10, 10, 11},
		"zero-plus": {0, 10, 1},
		// TODO: コインが1回分に満たない場合のテスト
		"plus-notenough": {1, 9, 1},
	}

	for name, tt := range cases {
		// TODO: ttをこのスコープで再定義しておく
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p := gacha.NewPlayer(tt.tickets, tt.coin)
			got := p.DrawableNum()
			if got != tt.want {
				// TODO: 分かりやすいメッセージを出してテストを失敗させる
				t.Errorf("want %d, but got %d", tt.want, got)
			}
		})
	}
}
