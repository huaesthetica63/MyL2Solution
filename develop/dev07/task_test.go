package main

import (
	"testing"
	"time"
)

//функция теста: проверяем полученное от функции время и сверяем его с тем, что выведет time.Now()
func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	chans := []<-chan interface{}{sig(2 * time.Second), sig(5 * time.Second), sig(1 * time.Second), sig(1 * time.Second), sig(10 * time.Second)}
	<-or(
		chans...,
	)
	//получили канал, закрытый только при условии, что все каналы из chans уже закрылись
	//теперь проверяем, так ли это получилось или нет
	//сюда запоминаем информацию по каждому каналу
	isClosedChan := []bool{false, false, false, false, false}
	_, ok := <-chans[0] //ок false если канал закрыт
	if !ok {
		isClosedChan[0] = true
	}
	_, ok = <-chans[1]
	if !ok {
		isClosedChan[1] = true
	}
	_, ok = <-chans[2]
	if !ok {
		isClosedChan[2] = true
	}
	_, ok = <-chans[3]
	if !ok {
		isClosedChan[3] = true
	}
	_, ok = <-chans[4]
	if !ok {
		isClosedChan[4] = true
	}
	for _, x := range isClosedChan {
		if !x { //если какой-то канал мы не увидели закрытым
			t.Errorf("Incorrect test\n") //тест провален
		}
	}

}
