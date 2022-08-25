package main

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
import (
	"fmt"
	"sort"
	"strings"
)

//функция для проверки двух слов - являются ли они анаграммами или нет
func CheckAnagramm(a, b string) bool {
	if len(a) != len(b) { //если они разной длины - уже все ясно
		return false //не анаграммы они
	}
	for _, x := range a { //обходим все буквы
		if strings.Count(a, string(x)) != strings.Count(b, string(x)) {
			//если количество одних и тех же букв не совпадает
			return false //не анаграммы
		}
	}
	return true //если все проверки пройдены, то слова являются анаграммами (состоят из одинакового набора букв)
}

//функция, выполняющая поставленную задачу
func AnagrammFind(arr []string) map[string][]string {
	res := make(map[string][]string) //карта с ключем-словом и значением - слайс слов
	for i, x := range arr {          //обходим весь массив, переводя слова  в нижний регистр
		arr[i] = strings.ToLower(x)
	}
	sort.Strings(arr) //сортируем массив , чтобы слова уже шли по порядку и мы не делали сортировок в карте
	for _, x := range arr {
		inMap := false          //проверяем, добавлена ли эта анаграмма в словарь или еще нет
		for k, v := range res { //обходим словарь
			if CheckAnagramm(k, x) { //проверяем все слова оттуда на анаграмму с текущим словом из массива
				res[k] = append(v, x) //если такая анаграмма уже есть, добавляем слово внутрь слайса
				inMap = true               //да, текущее слово в массиве уже записано в мапу
				break                      //прерываем цикл, потому что мы уже все выяснили
			}
		}
		if !inMap { //если анаграмма из слова не была найдена среди ключей
			res[x] = append(res[x], x) //делаем слово новым ключом и заодно вставляем его в слайс
		}
	}
	for k, v := range res { //теперь проверяем множества из одного слова и удаляем их
		if len(v) <= 1 {
			delete(res, k)
		}
	}
	return res
}
func main() {
	words := []string{"стОлиК", "пятАк", "пятка", "слиток", "тяпка", "листок", "слово"} //демонстрационный список слов
	fmt.Println(AnagrammFind(words)) //выводим получившийся словарь
}
