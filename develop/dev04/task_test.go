package main

import (
	"testing"
)

//сравнение мап из задания
func CompareMap(a, b map[string][]string) bool {
	for k, v := range a { //обходим первую мапу
		if len(v) != len(b[k]) { //сравниваем длины слайсов
			return false
		}
		for i, x := range b[k] { //сравниваем содержание слайсов поэлементно
			if v[i] != x {
				return false
			}
		}
	}
	return true
}

//функция тестирования, где мы проверим все примеры из задания
func TestAnnagramFind(t *testing.T) {
	got := AnagrammFind([]string{"стОлиК", "пятАк", "пятка", "слиток", "тяпка", "листок", "слово"})
	want := map[string][]string{"листок": {"листок", "слиток", "столик"}, "пятак": {"пятак", "пятка", "тяпка"}}
	if !CompareMap(got, want) {
		t.Errorf("got: %s want: %s\n", got, want)
	}
}
