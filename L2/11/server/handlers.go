package server

import "net/http"

var (
	CreateEventHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MOCKED"))
	}
	UpdateEventHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MOCKED"))
	}
	DeleteEventHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MOCKED"))
	}
	EventsForDayHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MOCKED"))
	}
	EventsForWeekHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MOCKED"))
	}
	EventsForMonthHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MOCKED"))
	}
)
