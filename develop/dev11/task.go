package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"slices"
	"strconv"
	"time"
)

// event - класс события календаря
type event struct {
	Date   time.Time `json:"date"`
	UserID int       `json:"user_id"`
}

// newEvent - конструктор класса event
func newEvent(date time.Time, userID int) *event {
	return &event{
		Date:   date,
		UserID: userID,
	}
}

// eventStorage - класс хранилища событий
type eventStorage struct {
	events []*event
}

// eventStorage - конструктор класса eventStorage
func newEventStorage() *eventStorage {
	return &eventStorage{
		events: make([]*event, 0),
	}
}

// addEvent - метод для добавления события в хранилище
func (es *eventStorage) addEvent(newEvent *event) {
	es.events = append(es.events, newEvent)
}

// removeEvent - метод для удаления события из хранилища
func (es *eventStorage) removeEvent(date time.Time, userID int) error {
	// поиск индекса удаляемого события
	eventIndex := slices.IndexFunc(es.events, func(e *event) bool {
		if e.Date == date && e.UserID == userID {
			return true
		}
		return false
	})
	// возврат ошибки, если событие не найдено
	if eventIndex == -1 {
		return errors.New("event not found")
	}
	slices.Delete(es.events, eventIndex, eventIndex+1) // удаление события
	return nil
}

// updateEvent - метод для обновления события
func (es *eventStorage) updateEvent(date time.Time, userID int, newDate time.Time) error {
	// поиск индекса обновляемого события
	eventIndex := slices.IndexFunc(es.events, func(e *event) bool {
		if e.Date == date && e.UserID == userID {
			return true
		}
		return false

	})
	// возврат ошибки, если событие не найдено
	if eventIndex == -1 {
		return errors.New("event not found")
	}
	es.events[eventIndex].Date = newDate // обновление события
	return nil
}

// findEventByDay - метод поиска событий по дню
func (es *eventStorage) findEventByDay(date time.Time) ([]*event, error) {
	// поиск событий и добавление их в слайс
	foundEvents := make([]*event, 0)
	for _, storedEvent := range es.events {
		if storedEvent.Date == date {
			foundEvents = append(foundEvents, storedEvent)
		}
	}
	// если слайс пустой, возврат ошибки
	if len(foundEvents) == 0 {
		return nil, errors.New("event not found")
	}
	return foundEvents, nil
}

// findEventByMonth - метод поиска событий по месяцу
func (es *eventStorage) findEventByMonth(date time.Time) ([]*event, error) {
	// поиск событий и добавление их в слайс
	foundEvents := make([]*event, 0)
	for _, storedEvent := range es.events {
		if storedEvent.Date.Month() == date.Month() && storedEvent.Date.Year() == date.Year() {
			foundEvents = append(foundEvents, storedEvent)
		}
	}
	// если слайс пустой, возврат ошибки
	if len(foundEvents) == 0 {
		return nil, errors.New("event not found")
	}
	return foundEvents, nil
}

// findEventByWeek - метод поиска событий по неделе
func (es *eventStorage) findEventByWeek(date time.Time) ([]*event, error) {
	// поиск событий и добавление их в слайс
	foundEvents := make([]*event, 0)
	for _, storedEvent := range es.events {
		_, storedWeek := storedEvent.Date.ISOWeek()
		_, searchWeek := date.ISOWeek()
		if storedWeek == searchWeek {
			foundEvents = append(foundEvents, storedEvent)
		}
	}
	// если слайс пустой, возврат ошибки
	if len(foundEvents) == 0 {
		return nil, errors.New("event not found")
	}
	return foundEvents, nil
}

// ключ от адреса сервера для контекста
var keyServerAddr = []byte("serverAddr")

// глобальное хранилище событий
var serverEventStorage = newEventStorage()

// класс ответа на неверный запрос
type badResponse struct {
	Error string `json:"error"`
}

// класс ответа на верный запрос
type goodResponse struct {
	Result string `json:"result"`
}

// createEvent - хэндлер функция для сервера, создающая новое событие
func createEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	hasDate := r.Form.Has("date")
	date := r.Form.Get("date")
	hasUserID := r.Form.Has("user_id")
	userID := r.Form.Get("user_id")

	convertedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	convertedUserID, err := strconv.Atoi(userID)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	serverEventStorage.addEvent(newEvent(convertedDate, convertedUserID))

	if !hasDate || !hasUserID {
		resp := badResponse{Error: "no event data was sent"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		resp := goodResponse{Result: "event was created"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// updateEvent - хэндлер функция для сервера, обновляющая событие
func updateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	hasDate := r.Form.Has("date")
	date := r.Form.Get("date")
	hasUserID := r.Form.Has("user_id")
	userID := r.Form.Get("user_id")
	hasNewDate := r.Form.Has("user_id")
	newDate := r.Form.Get("user_id")

	convertedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	convertedNewDate, err := time.Parse(time.DateOnly, newDate)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	convertedUserID, err := strconv.Atoi(userID)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	err = serverEventStorage.updateEvent(convertedDate, convertedUserID, convertedNewDate)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(503)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if !hasDate || !hasUserID || !hasNewDate {
		resp := badResponse{Error: "no event data was sent"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		resp := goodResponse{Result: "event was updated"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// deleteEvent - хэндлер функция для сервера, удаляющая событие
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	hasDate := r.Form.Has("date")
	date := r.Form.Get("date")
	hasUserID := r.Form.Has("user_id")
	userID := r.Form.Get("user_id")

	convertedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	convertedUserID, err := strconv.Atoi(userID)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	err = serverEventStorage.removeEvent(convertedDate, convertedUserID)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(503)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if !hasDate || !hasUserID {
		resp := badResponse{Error: "no event data was sent"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		resp := goodResponse{Result: "event was deleted"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// getEventsForDay - хэндлер функция для сервера, ищущая событие по дню
func getEventsForDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hasDate := r.URL.Query().Has("date")
	date := r.URL.Query().Get("date")

	convertedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	foundEvents, err := serverEventStorage.findEventByDay(convertedDate)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(503)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if !hasDate || foundEvents == nil {
		resp := badResponse{Error: "no event data was sent"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		responseText := ""
		for _, foundEvent := range foundEvents {
			responseText += "event: date=" + foundEvent.Date.String() + " user_id= \n"
		}
		resp := goodResponse{Result: responseText}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// getEventsForWeek - хэндлер функция для сервера, ищущая событие по неделе
func getEventsForWeek(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hasDate := r.URL.Query().Has("date")
	date := r.URL.Query().Get("date")

	convertedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	foundEvents, err := serverEventStorage.findEventByWeek(convertedDate)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(503)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if !hasDate || foundEvents == nil {
		resp := badResponse{Error: "no event data was sent"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		responseText := ""
		for _, foundEvent := range foundEvents {
			responseText += "event: date=" + foundEvent.Date.String() + " user_id= \n"
		}
		resp := goodResponse{Result: responseText}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// getEventsForMonth - хэндлер функция для сервера, ищущая событие по месяцу
func getEventsForMonth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	hasDate := r.URL.Query().Has("date")
	date := r.URL.Query().Get("date")

	convertedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	foundEvents, err := serverEventStorage.findEventByMonth(convertedDate)
	if err != nil {
		resp := badResponse{Error: err.Error()}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(503)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if !hasDate || foundEvents == nil {
		resp := badResponse{Error: "no event data was sent"}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		responseText := ""
		for _, foundEvent := range foundEvents {
			responseText += "event: date=" + foundEvent.Date.String() + " user_id= \n"
		}
		resp := goodResponse{Result: responseText}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {

	// создание сервера о объявление путей с привязкой хэндлер функций
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", createEvent)
	mux.HandleFunc("/update_event", updateEvent)
	mux.HandleFunc("/delete_event", deleteEvent)
	mux.HandleFunc("/events_for_day", getEventsForDay)
	mux.HandleFunc("/events_for_week", getEventsForWeek)
	mux.HandleFunc("/events_for_month", getEventsForMonth)

	// создание контекста
	ctx := context.Background()
	server := &http.Server{
		Addr:    ":4040",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	// запуск сервера
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
