package main

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

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
	"strconv"
	"strings"
)

var ( //флаги для использования утилиты (ключи)
	ColumnFlag  = flag.Int("k", 0, "sorting column number")
	NumericFlag = flag.Bool("n", false, "numeric sort")
	ReverseFlag = flag.Bool("r", false, "reverse order")
	UniqueFlag  = flag.Bool("u", false, "only unique values")
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

//для флага -r разворачиваем отсортированный слайс в обратном порядке
func ReverseLines(lines []string) {
	for i, j := 0, len(lines)-1; i < j; {
		lines[i], lines[j] = lines[j], lines[i]
		i++
		j--
	}
}

//создаем множество строк (set) без повторов (флаг -u)
func UniqueLines(lines []string) []string {
	helpmap := make(map[string]struct{}) //мапа с ключом-строкой без значения (пустая структура)
	for _, x := range lines {
		helpmap[x] = struct{}{}
	}
	var res []string //теперь записываем ключи из мапы в новый слайс
	for x := range helpmap {
		res = append(res, x)
	}
	return res
}

//сортировка числовых значений
func MySort(a []int) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

//флаг -n
func NumericSort(lines []string) ([]string, error) {
	var a []int               //слайс чисел для сортировки
	var s []string            //новый порядок строк
	for _, v := range lines { //обходим строки
		int, err := strconv.Atoi(strings.TrimSuffix(v, "\n")) //обрезаем символ \n и переводим в число
		if err != nil {
			return nil, errors.New("incorrect file")
		}
		a = append(a, int) //добавляем в слайс чисел
	}
	MySort(a)
	for _, v := range a {
		s = append(s, strconv.Itoa(v)+"\n") //записываем обратно в строки новый порядок чисел
	}
	lines = s

	return lines, nil
}

//структура для сортировки по колонкам
type columnWord struct {
	Word   string //слово из колонки
	ColNum int    //номер строки, откуда оно взято
}

//сортировка структур columnWord
func MySortColumn(a []columnWord) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i].Word > a[j].Word { //сортируем по слову
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

//сортировка по флагу -k
func ColumnSort(lines []string, colNum int) ([]string, error) {
	var sortSlice []columnWord //структура для хранения данных о каждой колонке
	for i, x := range lines {
		words := strings.Split(x, " ") //разделяем строку по пробелам
		if len(words) < colNum {       //если номер колонки больше числа слов, получившихся в строке
			return nil, errors.New("incorrect file")
		}
		sortSlice = append(sortSlice, columnWord{words[colNum-1], i}) //добавляем слово из нужной колонки и номер строки
	}
	MySortColumn(sortSlice) //сортировка
	var res []string        //переписываем строки в новом порядке в соответствии со структурой columnWord
	for _, x := range sortSlice {
		res = append(res, lines[x.ColNum])
	}
	return res, nil
}

//запись слайса строк в указанный файл
func WriteFile(filename string, lines []string) error {
	outputfile, err := os.Create(filename)
	if err != nil {
		return errors.New("incorrect file")
	}
	output := bufio.NewWriter(outputfile)
	defer outputfile.Close()
	defer output.Flush() //очищаем файловый буфер  в конце
	for _, str := range lines {
		_, err := output.WriteString(str)
		if err != nil {
			return errors.New("incorrect file")
		}
	}
	return nil
}

//простая сортировка без флагов по строкам
func SimpleSort(a []string) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}
func main() {
	flag.Parse()                           //считываем флаги с команды терминала
	iFileName := os.Args[len(os.Args)-2]   //входной файл
	outFileName := os.Args[len(os.Args)-1] //файл для результата
	lines, err := ReadFile(iFileName)
	if err != nil {
		log.Fatal(err)
	}
	if *NumericFlag { //сортировка для чисел
		lines, err = NumericSort(lines) //числовая сортировка
		if err != nil {
			log.Fatal(err)
		}
	} else if *ColumnFlag > 0 { //сортировка по колонкам
		lines, err = ColumnSort(lines, *ColumnFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		SimpleSort(lines)
	}
	if *ReverseFlag {
		ReverseLines(lines)
	}
	if *UniqueFlag {
		UniqueLines(lines)
	}
	fmt.Println(lines)
	err = WriteFile(outFileName, lines) //записываем в файл результат
	if err != nil {
		log.Fatal(err)
	}
}
