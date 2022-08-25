package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
//наша функция по распаковке
func UnpackString(str string) (string, error) {
	var res strings.Builder            //билдер строки, так конкатенация рун будет быстрее, чем обычный "+"
	strarr := []rune(str)              //получаем слайс рун, чтобы обрабатывать посимвольно, а не побайтово
	var curRune *rune                  //текущий символ распаковки
	curRune = nil                      //пока у нас нет символа, который мы запомнили для записи
	for i := 0; i < len(strarr); i++ { //обходим посимвольно слайс рун
		currchar := strarr[i] //получаем текущий символ
		if currchar == '\\' { //если это символ экранирования
			if i+1 >= len(strarr) { //если \ - последний символ в строке
				return "", errors.New("incorrect string")
			}
			i++                 //пропускаем его и записываем следующий за ним как "смысловой" символ, который может быть вписан в строку
			if curRune != nil { // если у нас уже есть запомненный символ
				res.WriteRune(*curRune) //прописываем его один раз
			} else {
				curRune = new(rune) //выделяем память под руну
			}
			*curRune = strarr[i] //запоминаем новый символ после символа экранирования
			continue             //переходим к следующей итерации, не проходя проверки ниже
		}
		if !unicode.IsDigit(currchar) { //если попался символ
			if curRune != nil { // если у нас уже есть запомненный символ
				res.WriteRune(*curRune) //прописываем его один раз
			} else {
				curRune = new(rune) //выделяем память под руну
			}
			*curRune = currchar // запоминаем текущий символ
		} else { //если попалась цифра
			count, _ := strconv.Atoi(string(currchar)) //переводим в число количество символов на распаковку
			if curRune == nil {                        //если нет символа для записи
				return "", errors.New("incorrect string") //это неправильная строка (как в случае с "45")
			}
			for j := 0; j < count; j++ { //распаковываем
				res.WriteRune(*curRune)
			}
			curRune = nil //снова обнуляем символ
		}
	}
	if curRune != nil { //если остался символ, который мы запомнили ("abcd" - d будет таким символом)
		res.WriteRune(*curRune) //прописываем его один раз
	}
	return res.String(), nil
}
func main() {
	str, err := UnpackString(`Yo9, \73 \\\7`)
	if err == nil {
		fmt.Println(str)
	} else {
		fmt.Println("Ошибка!")
		log.Fatal(err)
	}
}
