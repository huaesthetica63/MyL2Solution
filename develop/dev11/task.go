package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Logger struct { //обертка для хендлера
	handler http.Handler
}

//такая обертка нужна для переопределения метода ServeHTTP
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//от стандартного он отличается логом
	initTime := time.Now()
	l.handler.ServeHTTP(w, r)
	endTime := time.Now()
	diff := endTime.Sub(initTime)
	log.Printf("%s %s %v", r.Method, r.URL.Path, diff) //выводим время обработки хендлера
}
func WrapHandler(h http.Handler) *Logger { //оборачиваем стандартный handler
	var wrapLog Logger
	wrapLog.handler = h
	return &wrapLog
}

type Event struct { //события, которые используются в работе календаря
	EventID int       `json:"event_id"` //уникальный идентификатор
	Content string    `json: "content"` //содержание (описание для этой даты)
	Date    time.Time `json:"date"`     //дата
}

//создаем глобально мьютекс для управления процессами, а также набор событий, мьютекс обеспечит безопасный доступ к нему
var mutex *sync.Mutex
var events []Event //слайс событий
func ParseJSON(r *http.Request) (Event, error) { //функция для парсинга ответа
	var ev Event
	err := json.NewDecoder(r.Body).Decode(&ev) //пробуем парсить json в структуру
	if err != nil {                            //если не получилось
		return ev, errors.New("incorrect json") //возвращаем ошибку
	}
	return ev, nil
}
func ValidateEvent(ev Event) error { //валидация event с точки зрения логики
	if ev.EventID <= 0 || ev.Content == "" { //проверяем поля на логику
		return errors.New("invalid event")
	}
	return nil //это если все ок и данные валидны
}
func CreateNewEvent(ev Event) error { //функция для создания нового события
	mutex.Lock() //блокируем остальным процессам доступ к данным
	defer mutex.Unlock()
	//проверим, есть ли уже такой id или нет
	for _, x := range events {
		if x.EventID == ev.EventID {
			return errors.New("eventId is already exist")
		}
	}
	events = append(events, ev) //добавляем событие в слайс событий
	return nil
}
func CreateEvent(w http.ResponseWriter, r *http.Request) { //обработчик для create_event
	if r.Method != http.MethodPost { //если это не post-метод, дальше не обрабатываем его
		errorResponse(w, "not POST-method", http.StatusBadRequest)
		return
	}
	newEvent, err := ParseJSON(r) //парсим json
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ValidateEvent(newEvent)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	//прошли парсинг и валидацию
	if err := CreateNewEvent(newEvent); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultResponse(w, "successfully!", []Event{newEvent}, http.StatusCreated)
}

//обработчик обновления события
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //если это не post-метод, дальше не обрабатываем его
		errorResponse(w, "not POST-method", http.StatusBadRequest)
		return
	}
	newEvent, err := ParseJSON(r) //парсим json
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ValidateEvent(newEvent)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	//прошли парсинг и валидацию
	mutex.Lock() //пока меняем слайс, другие потоки не должны туда влезать
	defer mutex.Unlock()
	for i, x := range events {
		if x.EventID == newEvent.EventID {
			events[i] = newEvent
			return
		}
	}
	return
}
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //если это не post-метод, дальше не обрабатываем его
		errorResponse(w, "not POST-method", http.StatusBadRequest)
		return
	}
	newEvent, err := ParseJSON(r) //парсим json
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ValidateEvent(newEvent)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	//прошли парсинг и валидацию
	mutex.Lock() //пока меняем слайс, другие потоки не должны туда влезать
	defer mutex.Unlock()
	for i, x := range events {
		if x.EventID == newEvent.EventID {
			events = append(events[:i], events[i+1:]...)
		}
	}
	return
}
func EventsByDay(date time.Time) []Event { //вернем список с полным совпадением дат вплоть до дня
	var res []Event
	mutex.Lock()
	mutex.Unlock()
	for _, x := range events {
		if x.Date.Year() == date.Year() && x.Date.Month() == date.Month() && x.Date.Day() == date.Day() {
			res = append(res, x)
		}
	}
	return res
}
func EventsByWeek(date time.Time) []Event { //вернем список с совпадением до недели
	var res []Event
	mutex.Lock()
	mutex.Unlock()
	for _, x := range events {
		dif := date.Sub(x.Date)
		if dif < 0 {
			dif = -dif
		}
		if dif <= time.Duration(7*24)*time.Hour {
			res = append(res, x)
		}
	}
	return res
}
func EventsByMonth(date time.Time) []Event { //вернем список с совпадением дат вплоть до месяца
	var res []Event
	mutex.Lock()
	mutex.Unlock()
	for _, x := range events {
		if x.Date.Year() == date.Year() && x.Date.Month() == date.Month() {
			res = append(res, x)
		}
	}
	return res
}
func EventsDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { //если это не get-метод, дальше не обрабатываем его
		errorResponse(w, "not Get-method", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2019-09-09", r.URL.Query().Get("date")) //получаем дату и парсим по шаблону
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := EventsByDay(date)
	resultResponse(w, "successfully", res, http.StatusOK)
}
func EventsWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { //если это не get-метод, дальше не обрабатываем его
		errorResponse(w, "not Get-method", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2019-09-09", r.URL.Query().Get("date")) //получаем дату и парсим по шаблону
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := EventsByWeek(date)
	resultResponse(w, "successfully", res, http.StatusOK)
}
func EventsMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { //если это не get-метод, дальше не обрабатываем его
		errorResponse(w, "not Get-method", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2019-09-09", r.URL.Query().Get("date")) //получаем дату и парсим по шаблону
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := EventsByMonth(date)
	resultResponse(w, "successfully", res, http.StatusOK)
}

func main() {
	port := os.Getenv("PORT")             //читаем порт из переменных окружения
	httpMultiplexor := http.NewServeMux() //мультиплексор, который выбирает обработчик в зависимости от запроса
	//объявляем обработчики
	httpMultiplexor.HandleFunc("/create_event", CreateEvent)
	httpMultiplexor.HandleFunc("/update_event", UpdateEvent)
	httpMultiplexor.HandleFunc("/delete_event", DeleteEvent)
	httpMultiplexor.HandleFunc("/events_for_day", EventsDay)
	httpMultiplexor.HandleFunc("/events_for_week", EventsWeek)
	httpMultiplexor.HandleFunc("/events_for_month", EventsMonth)
	// middleware-логгер
	midLogger := WrapHandler(httpMultiplexor)         //по сути мы создали обертку для мультиплексора
	log.Fatalln(http.ListenAndServe(port, midLogger)) //запускаем сервер
}

func errorResponse(w http.ResponseWriter, e string, status int) {
	errorResponse := struct {
		Error string `json:"error"`
	}{Error: e}

	js, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func resultResponse(w http.ResponseWriter, r string, e []Event, status int) {
	resultResponse := struct {
		Result string  `json:"result"`
		Events []Event `json:"events"`
	}{Result: r, Events: e}

	js, err := json.Marshal(resultResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
