package main

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	FieldsFlag    = flag.Int("f", 0, "fields")
	DelimiterFlag = flag.String("d", "\t", "delimiter")
	SeparatedFlag = flag.Bool("s", false, "separated")
)

func main() {
	flag.Parse()                          //получаем флаги с консоли
	scanner := bufio.NewScanner(os.Stdin) //чтение со стандартного потока ввода
	fmt.Println(`"exit" - для завершения ввода`)
	var words [][]string //слайс слайсов строк (каждую строку дробим по разделителям)
	for {
		ok := scanner.Scan()
		if !ok {
			log.Fatal(errors.New("ошибка чтения с потока"))
		}
		line := scanner.Text() //получаем текст - одну строку
		if line == "exit" {    //если эта строка exit
			break //ввод завершен
		}
		if !(*SeparatedFlag && !strings.Contains(line, *DelimiterFlag)) {
			words = append(words, strings.Split(line, *DelimiterFlag)) //если выставлен флаг и в строке нет разделителя
			//игнорируем ее, если все ок - записываем ее в words
		}

	}
	//проверяем: выставлен ли флаг по полям
	if *FieldsFlag < 0 { //отрицательные значения некорректны
		log.Fatal(errors.New("некорректный флаг консоли"))
	}
	if *FieldsFlag != 0 { //если поле ненулевое
		var columns []string
		for _, slice := range words {
			columns = append(columns, slice[*FieldsFlag])
		}
		fmt.Println(columns)
	} else { //если флаг не выставлен, печатаем как есть
		for _, slice := range words { //сначала обходим послайсово
			for _, word := range slice { //внутри каждого слайса обходим слова
				fmt.Print(word + *DelimiterFlag) //печатаем с разделителем
			}
			fmt.Println("") // перевод строки
		}
	}
}
