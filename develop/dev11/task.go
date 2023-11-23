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

type event struct {
	Date   time.Time `json:"date"`
	UserId int       `json:"user_id"`
}

func newEvent(date time.Time, userId int) *event {
	return &event{
		Date:   date,
		UserId: userId,
	}
}

type eventStorage struct {
	events []*event
}

func newEventStorage() *eventStorage {
	return &eventStorage{
		events: make([]*event, 0),
	}
}

func (es *eventStorage) addEvent(newEvent *event) {
	es.events = append(es.events, newEvent)
}

func (es *eventStorage) removeEvent(date time.Time, userId int) error {
	eventIndex := slices.IndexFunc(es.events, func(e *event) bool {
		if e.Date == date && e.UserId == userId {
			return true
		} else {
			return false
		}
	})
	if eventIndex == -1 {
		return errors.New("event not found")
	}
	slices.Delete(es.events, eventIndex, eventIndex+1)
	return nil
}

func (es *eventStorage) updateEvent(date time.Time, userId int, newDate time.Time) error {
	eventIndex := slices.IndexFunc(es.events, func(e *event) bool {
		if e.Date == date && e.UserId == userId {
			return true
		} else {
			return false
		}
	})
	if eventIndex == -1 {
		return errors.New("event not found")
	}
	es.events[eventIndex].Date = newDate
	return nil
}

func (es *eventStorage) findEventByDay(date time.Time) ([]*event, error) {
	foundEvents := make([]*event, 0)
	for _, storedEvent := range es.events {
		if storedEvent.Date == date {
			foundEvents = append(foundEvents, storedEvent)
		}
	}
	if len(foundEvents) == 0 {
		return nil, errors.New("event not found")
	}
	return foundEvents, nil
}

func (es *eventStorage) findEventByMonth(date time.Time) ([]*event, error) {
	foundEvents := make([]*event, 0)
	for _, storedEvent := range es.events {
		if storedEvent.Date.Month() == date.Month() && storedEvent.Date.Year() == date.Year() {
			foundEvents = append(foundEvents, storedEvent)
		}
	}
	if len(foundEvents) == 0 {
		return nil, errors.New("event not found")
	}
	return foundEvents, nil
}

func (es *eventStorage) findEventByWeek(date time.Time) ([]*event, error) {
	foundEvents := make([]*event, 0)
	for _, storedEvent := range es.events {
		_, storedWeek := storedEvent.Date.ISOWeek()
		_, searchWeek := date.ISOWeek()
		if storedWeek == searchWeek {
			foundEvents = append(foundEvents, storedEvent)
		}
	}
	if len(foundEvents) == 0 {
		return nil, errors.New("event not found")
	}
	return foundEvents, nil
}

var keyServerAddr = "serverAddr"

var serverEventStorage = newEventStorage()

type badResponse struct {
	Error string `json:"error"`
}

type goodResponse struct {
	Result string `json:"result"`
}

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
	hasUserId := r.Form.Has("user_id")
	userId := r.Form.Get("user_id")

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

	convertedUserId, err := strconv.Atoi(userId)
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

	serverEventStorage.addEvent(newEvent(convertedDate, convertedUserId))

	if !hasDate || !hasUserId {
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
	hasUserId := r.Form.Has("user_id")
	userId := r.Form.Get("user_id")
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

	convertedUserId, err := strconv.Atoi(userId)
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

	err = serverEventStorage.updateEvent(convertedDate, convertedUserId, convertedNewDate)
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

	if !hasDate || !hasUserId || !hasNewDate {
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
	hasUserId := r.Form.Has("user_id")
	userId := r.Form.Get("user_id")

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

	convertedUserId, err := strconv.Atoi(userId)
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

	err = serverEventStorage.removeEvent(convertedDate, convertedUserId)
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

	if !hasDate || !hasUserId {
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

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", createEvent)
	mux.HandleFunc("/update_event", updateEvent)
	mux.HandleFunc("/delete_event", deleteEvent)
	mux.HandleFunc("/events_for_day", getEventsForDay)
	mux.HandleFunc("/events_for_week", getEventsForWeek)
	mux.HandleFunc("/events_for_month", getEventsForMonth)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
