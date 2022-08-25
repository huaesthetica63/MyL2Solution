package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Паттерн "цепочка вызовов" представляет собой цепочку обработчиков, которая принимает на вход
некий запрос и передает его обработчика к обработчику. Каждый из них решает: сможет ли он
обработать этот запрос или нужно передать его по цепочке далее (паспортный стол или регистратура - пример паттерна из жизни)
Паттерн применяется в том случае, когда система должна обрабатывать разнообразные запросы, суть которых и методы
обработки заранее неизвестны, когда обработчики могут меняться динамически (например, при мастштабировании системы),
или когда их порядок должен быть строго последовательным (например, обработчики можно расположить по степени
сложности запроса - от самого простого случая, к самому трудному - первый обработчик обрабатывает самые рядовые
базовые вещи, последний в цепочке - редкие трудные случаи)

Реальный пример использования: допустим, наше приложение подключается к удаленному серверу, либо он вернет
нам ожидаемый результат (валидный json или что-то такое), либо же код ошибки - эти данные мы передаем цепочке,
самый первый обработчик принимает валидные данные, если они не валидны, он передает их дальше по обработчикам,
которые занимаются обработкой конкретных ошибок

Плюсы: применение паттерна позволяет уменьшить зависимость между клиентом и обработчиком запроса,
количество и структуру обработчиков можно легко поменять, сами обработчики как элементы общей цепочки станут
более единообразными, похожими друг на друга, что улучшит читабельность кода

Минусы: обработчики могут столкнуться с таким запросом, который не сможет обработать ни один элемент цепи,
поэтому есть риск того, что запрос может пройти через всю цепочку безрезультатно

*/

import (
	"encoding/json"
	"fmt"
)

type Response struct { //то, что получили от сервера
	message string //само сообщение (валидный json или код ошибки)
}
type ChainElement interface { //интерфейс, описывающий поведение одного элемента цепи (обработчика)
	Execute(res *Response)   // обработать запрос
	SetNext(ch ChainElement) //перейти к следующему обработчику
}
type ValidMessagePrinter struct { //обработчик на случай получения валидной строки json-формата
	Next ChainElement //следующий обработчик
}

func (vmp *ValidMessagePrinter) Execute(res *Response) { //обрабатываем ответ сервера , предполагая, что это json
	var jsonUnmarshal interface{}                              //пустой интерфейс используется из-за незнания структуры json'а - она может быть любой
	err := json.Unmarshal([]byte(res.message), &jsonUnmarshal) //попробуем расшифровать строку как json
	if err != nil {
		fmt.Println("Невалидный JSON!")
		vmp.Next.Execute(res) //передали следующему обработчику, пусть сам разбирается
	} else {
		fmt.Println("С сервера получен JSON: ", res.message)
	}
}
func (vmp *ValidMessagePrinter) SetNext(ch ChainElement) {
	vmp.Next = ch
}

type Error404Printer struct { //обработчик на случай ошибки 404
	Next ChainElement //следующий обработчик
}

func (err404 *Error404Printer) Execute(res *Response) { //обрабатываем ответ сервера только если это 404 ошибка
	if res.message == "404 error" {
		fmt.Println("Ошибка 404!") // обработчик знает такую ошибку
	} else {
		fmt.Println("С сервера получена не 404 ошибка... ")
		err404.Next.Execute(res)
	}
}
func (err404 *Error404Printer) SetNext(ch ChainElement) {
	err404.Next = ch
}

type UnknownErrorPrinter struct { //обработчик на случай,если это не json, и не 404
	Next ChainElement //следующий обработчик
}

func (uerr *UnknownErrorPrinter) Execute(res *Response) { //обрабатываем ответ сервера
	var jsonUnmarshal interface{}                              //пустой интерфейс используется из-за незнания структуры json'а - она может быть любой
	err := json.Unmarshal([]byte(res.message), &jsonUnmarshal) //попробуем расшифровать строку как json
	if (err == nil) || (res.message == "404 error") {          //если это валидный json или 404
		uerr.Next.Execute(res) //мы с такими строками тут не работаем
	} else {
		fmt.Println("С сервера получена неизвестная ошибка!")
	}
}
func (uerr *UnknownErrorPrinter) SetNext(ch ChainElement) {
	uerr.Next = ch
}
func main() {
	//создаем цепочку обработчиков
	var jsonHandler = &ValidMessagePrinter{}
	var err404Handler = &Error404Printer{}
	var unknHandler = &UnknownErrorPrinter{}
	jsonHandler.Next = err404Handler
	err404Handler.Next = unknHandler
	var str1 = `{"a":"b"}` //json
	var resp Response
	resp.message = str1
	jsonHandler.Execute(&resp)
	var str2 = "404 error" //не json
	resp.message = str2
	jsonHandler.Execute(&resp)
	var str3 = "пупа и лупа получили зарплату.... " //неизвестный ответ сервера
	resp.message = str3
	jsonHandler.Execute(&resp)
}
