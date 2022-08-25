package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func wget(url string) error {
	response, err := http.Get(url) //получаем ответ по протоколу http по данному url
	if err != nil {
		return errors.New("ошибка по данному url")
	}
	temp := strings.Split(url, "/")      //получаем составные части url для выяснения имя файла
	filename := temp[len(temp)-1]        //последняя часть - имя файла
	saveFile, err := os.Create(filename) //создаем файл на компьютере для сохранения файла по url
	if err != nil {
		return errors.New("ошибка создания файла")
	}
	defer saveFile.Close() //закрываем в конце файл
	_, err = io.Copy(saveFile, response.Body)
	if err != nil {
		return errors.New("ошибка при сохранении файла")
	}
	fmt.Println("Файл успешно сохранен!")
	return nil
}
func main() {
	scanner := bufio.NewScanner(os.Stdin) //чтение со стандартного потока ввода
	fmt.Println(`Введите url: `)
	ok := scanner.Scan()
	if !ok {
		log.Fatal("ошибка")
	}
	url := scanner.Text()
	err := wget(url)
	if err != nil {
		log.Fatal(err)
	}
}
