package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/
//наша функция, представляющая or-канал: совокупность всех каналов из слайса channels,
//канал закрывается, когда закрывается каждый из этих каналов
func or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 { //проверяем, есть ли вообще элементы в слайсе
		return nil
	}
	if len(channels) == 1 { //если в слайсе один канал - его же и возвращаем
		return channels[0]
	}
	reschan := make(chan interface{}) //or-канал
	go func() {                       //выводим остальную часть в рамки горутины, потому что каналы могут закрываться очень долго,
		//поэтому мы сразу вернем or-канал, а закроем его только по завершении этой горутины
		defer close(reschan)  //закрываем его в конце горутины
		<-channels[0]         //ждем завершения первого в слайсе канала
		<-or(channels[1:]...) //прибегаем  к использованию рекурсии, передавая в новый вызов слайс без первого элемента
	}()
	return reschan //возвращаем канал
}
func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	//сделаем чуть-чуть поменьше времени на каждый канал, потому что ждать минуту или даже час - слишком долго
	start := time.Now()
	<-or(
		sig(2*time.Second),
		sig(5*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(10*time.Second),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
