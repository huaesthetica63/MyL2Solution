package main

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

//функция обработки команд
func ExecComm(comm []string) error {
	for _, x := range comm { //обходим список команд
		args := strings.Split(x, " ") //каждая команда может содержать список аргументов через пробел
		commandName := args[0]        //первое слово в строке - имя самой команды
		switch commandName {
		case "cd":
			if len(args) < 2 { //если agrs меньше 2 элементов, значит, мы не передали аргументы
				return errors.New("недостаточно аргументов")
			}
			dir := args[1] //новая директория
			fmt.Println("Смена директории на %s", dir)
			os.Chdir(dir) //сменили текущую директорию
		case "pwd":
			dir, _ := os.Getwd()
			fmt.Println("Текущая директория: ", dir) //вывели директорию
		case "echo":
			if len(args) < 2 {
				return errors.New("недостаточно аргументов")
			}
			for i := 0; i < len(args); i++ {
				fmt.Fprintf(os.Stdout, args[i+1]+" ") //записываем аргументы в поток вывода
			}
		case "kill":
			if len(args) < 2 {
				return errors.New("недостаточно аргументов")
			}
			err := exec.Command("kill", args[1]).Run()
			if err != nil {
				return err
			}
		case "ps":
			procs, err := ps.Processes()
			if err != nil {
				return err
			}
			fmt.Println("Процессы: ")
			for _, v := range procs {
				fmt.Println(v.Pid())
			}
		default: //неизвестная команда, невошедшая в список выше
			return errors.New("неизвестная команда")
		}
	}
	return nil
}
func main() {
	scanner := bufio.NewScanner(os.Stdin) //чтение со стандартного потока ввода
	fmt.Println(`Для выхода введите команду "exit"`)
	for { //бесконечный цикл с "прослушкой" команд терминала
		curDir, _ := os.Getwd()    //получаем текущую директорию
		fmt.Printf("%s> ", curDir) //выводим предложение к вводу команд с указанием текущей директории
		ok := scanner.Scan()
		if !ok {
			log.Fatal(errors.New("ошибка при вводе команды"))
		}
		enter := scanner.Text() //читаем ввод пользователя
		if enter == "exit" {    //если ввели строку exit - выходим из всей программы
			fmt.Println("Выход...")
			os.Exit(0)
		}
		comms := strings.Split(enter, "|") // разделитель для использования множественных команд в одной строке
		err := ExecComm(comms)
		if err != nil {
			log.Fatal(err)
		}
	}
}
