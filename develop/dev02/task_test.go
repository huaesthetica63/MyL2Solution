package main

import (
	"testing"
)

//функция тестирования, где мы проверим все примеры из задания
func TestUnpackString(t *testing.T) {
	inputArr := []string{`a4bc2d5e`, `abcd`, ``, `qwe\4\5`, `qwe\45`, `qwe\\5`}
	wantArr := []string{`aaaabccddddde`, `abcd`, ``, `qwe45`, `qwe44444`, `qwe\\\\\`}
	for i, x := range inputArr {
		got, err := UnpackString(x)
		if err != nil {
			t.Errorf("got incorrect string\n")
		}
		want := wantArr[i]
		if got != want {
			t.Errorf("got: %s want: %s\n", got, want)
		}
	}

}
