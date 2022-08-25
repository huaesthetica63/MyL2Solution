package main

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

var (
	AfterFlag   = flag.Int("A", 0, "after")
	BeforeFlag  = flag.Int("B", 0, "before")
	ContextFlag = flag.Int("C", 0, "context")
	CountFlag   = flag.Bool("c", false, "count")
	IgnoreFlag  = flag.Bool("i", false, "ignore-case")
	InvertFlag  = flag.Bool("v", false, "invert")
	FixedFlag   = flag.Bool("F", false, "fixed")
	LineFlag    = flag.Bool("n", false, "line num")
)

//считываем файл построчно и выдаем слайс из строк
func ReadFile(filename string) ([]string, error) {
	inputfile, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("incorrect file")
	}
	defer inputfile.Close()
	input := bufio.NewReader(inputfile) //ридер с файла
	var res []string
	for {
		str, err := input.ReadString('\n') //разделитель - перевод строки
		if err != nil {
			if err == io.EOF { //если дошли до конца строки - выходим из цикла
				res = append(res, str+"\n")
				break
			}
			return nil, errors.New("incorrect file") //во всех остальных случаях - это ошибка файла
		}
		res = append(res, str) //добавляем строчку к слайсу
	}
	return res, nil
}

type StrGrep struct { //структура для хранения результата grep
	text string //содержание строки
	line int    //ее номер
}

func FindGrep(lines []string, pattern string) []StrGrep {
	var grep *regexp.Regexp //регулярное выражение для поиска
	if *ContextFlag > 0 {   //если контекст включен - включаем after и before одновременно
		*AfterFlag = *ContextFlag
		*BeforeFlag = *ContextFlag
	}
	var res []StrGrep
	if *IgnoreFlag {
		grep = regexp.MustCompile("(?i)" + pattern) //для игнорирования регистра нужно добавить модификатор режима (?i)
	} else {
		grep = regexp.MustCompile(pattern)
	}
	for i, v := range lines { //обходим строки
		if *FixedFlag { //сравниваем строки как обычные строки, а не как регулярные выражения
			if v == pattern+"\n" {
				res = append(res, StrGrep{v, i + 1}) //заносим строку и ее порядковый номер
			}
		} else {
			str := grep.FindString(v)      //если флаг fixed не стоит, обрабатываем строки как регулярные выражения
			if str == "" && !*InvertFlag { //если это не строка с рег.выражением и у нас отключен флаг
				continue //переходим к следующей итерации
			} else if (str == "" && *InvertFlag) || (str != "" && !*InvertFlag) {
				//при остальных условиях - записываем эту строку и проверяем флаги
				if *BeforeFlag > 0 { //если надо записать строки до
					startIndex := i - *BeforeFlag //индекс, откуда начинать
					if startIndex < 0 {           //если N строк не получается
						startIndex = 0 //начинаем записывать с первой и до исходной
					}
					for ; startIndex < i; startIndex++ {
						res = append(res, StrGrep{lines[startIndex], startIndex + 1}) //заносим строку и ее порядковый номер
					}
				}
				res = append(res, StrGrep{v, i + 1})
				if *AfterFlag > 0 { //алгоритм аналогичен
					for startIndex := i + 1; startIndex < len(lines); startIndex++ {
						if startIndex-i > *AfterFlag { //если уже записали N строк
							break //все
						}
						res = append(res, StrGrep{lines[startIndex], startIndex + 1})
					}

				}
			}
		}
	}
	return res
}
func main() {
	flag.Parse()
	inputFile := os.Args[0]
	pattern := os.Args[1]
	lines, err := ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	res := FindGrep(lines, pattern)
	if *CountFlag { //если включен флаг count
		fmt.Println(len(res)) //выводим количество строк
	} else {
		for _, x := range res {

			if *LineFlag {
				fmt.Printf("%d) %s", x.line, x.text)
			} else {
				fmt.Printf("%s", x.text)
			}
		}

	}

}
