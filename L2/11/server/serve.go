package server

import (
	"net/http"
	"strings"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	p := r.URL.Path
	switch {
	case strings.HasPrefix(p, "/create_event"):
		h = postJsonHandlerBuilder(CreateEventHandler)
	case strings.HasPrefix(p, "/update_event"):
		h = postJsonHandlerBuilder(UpdateEventHandler)
	case strings.HasPrefix(p, "/delete_event"):
		h = postJsonHandlerBuilder(DeleteEventHandler)
	case strings.HasPrefix(p, "/events_for_day"):
		h = getJsonHandlerBuilder(EventsForDayHandler)
	case strings.HasPrefix(p, "/events_for_week"):
		h = getJsonHandlerBuilder(EventsForWeekHandler)
	case strings.HasPrefix(p, "/events_for_month"):
		h = getJsonHandlerBuilder(EventsForMonthHandler)
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}
