package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeoutFlag := flag.Int("timeout", 10, "timeout")
	flag.Parse()
	exitSig := make(chan os.Signal)
	if os.Args[0] != "go-telnet" {
		os.Exit(1)
	}
	timeout := time.Duration(*timeoutFlag) * time.Second //по умолчанию 10 секунд
	host := os.Args[len(os.Args)-2]
	port := os.Args[len(os.Args)-1]
	signal.Notify(exitSig, os.Interrupt, os.Kill)   //ctrl+d это зарезервированный сигнал kill
	connectionParam := net.JoinHostPort(host, port) //сокет с таймаутом
	conn, err := net.DialTimeout("tcp", connectionParam, timeout)
	if err != nil {
		log.Fatal("ошибка подключения")
	}
	go func() { //прослушиваем на получение сигнала kill (ctrl+d)
		<-exitSig //считываем сигнал
		fmt.Println("\nВыход...")
		conn.Close() //закрываем сокет
		os.Exit(0)   //завершаем программу
	}()
	//прослушиваем на получение сообщений с сервера
	go func() {
		_, err := io.Copy(conn, os.Stdout)
		if err != nil {
			log.Fatal("ошибка с сервера")
		} else {
			fmt.Println("Сообщение отправлено!")
		}
	}()
	//прослушиваем на отсылку сообщений с stdin
	go func() {
		_, err := io.Copy(os.Stdin, conn)
		if err != nil {
			log.Fatal("ошибка во время отправки сообщения")
		}
	}()
}
