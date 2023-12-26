package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"httprestapi/model"
	"httprestapi/server"
	"net/http"
	"time"
)

var successmessage = `{"result":"success"}`
var badinputmessage = `{"error":"bad input"}`
var serviceunavailablemessage = `{"error":"service unavailable"}`

func init() {
	calendar := model.NewCalendar()
	server.CreateEventHandler = buildEventHandler(calendar.AddEvent)
	server.UpdateEventHandler = buildEventHandler(calendar.UpdateEvent)
	server.DeleteEventHandler = buildEventIdHandler(calendar.DeleteEvent)
	server.EventsForDayHandler = buildCalendarHandler(calendar.GetEventsForDay)
	server.EventsForWeekHandler = buildCalendarHandler(calendar.GetEventsForWeek)
	server.EventsForMonthHandler = buildCalendarHandler(calendar.GetEventsForMonth)
}

// функция билдер обработчик евент айди
func buildEventIdHandler(f func(id string) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.Form.Get("id")
		err := f(id)
		if err != nil {
			errorNoHeaders(w, serviceunavailablemessage, http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte(successmessage))
	}
}

// функция билдер обработчик евентов
func buildEventHandler(f func(e model.Event) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := parseRequestForm(r)
		if err != nil {
			errorNoHeaders(w, badinputmessage, http.StatusBadRequest)
			return
		}
		err = f(event)
		if err != nil {
			errorNoHeaders(w, serviceunavailablemessage, http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte(successmessage))
	}
}

// функция билдер запросов на интервалы
func buildCalendarHandler(f func(t time.Time) ([]model.Event, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := parseUrlQuery(r)
		if err != nil {
			errorNoHeaders(w, serviceunavailablemessage, http.StatusServiceUnavailable)
			return
		}
		events, err := f(date)
		if err != nil {
			errorNoHeaders(w, serviceunavailablemessage, http.StatusServiceUnavailable)
			return
		}
		message, err := json.Marshal(&events)
		if err != nil {
			errorNoHeaders(w, serviceunavailablemessage, http.StatusServiceUnavailable)
			return
		}
		w.Write(message)
	}
}

// тоже самое что и стандартный http.Error но без ненужных Headers
func errorNoHeaders(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

// хелпер функция для парсигна формы
func parseRequestForm(r *http.Request) (model.Event, error) {
	r.ParseForm()
	id := r.Form.Get("id")
	descr := r.Form.Get("description")
	sdate := r.Form.Get("date")
	date, err := parseDate(sdate)
	if err != nil {
		return model.Event{}, errors.New("parce form error")
	}
	return model.Event{
		Id:          id,
		Description: descr,
		Date:        date,
	}, nil
}

// хелпер функция для парсинга url query
func parseUrlQuery(r *http.Request) (time.Time, error) {
	sdate := r.URL.Query().Get("date")
	return parseDate(sdate)
}
func parseDate(sdate string) (time.Time, error) {
	date, err := time.Parse("2006-01-02 15:04", sdate)
	if err != nil {
		date, err = time.Parse("2006-01-02", sdate)
		if err != nil {
			return date, err
		}
	}
	return date, nil
}
