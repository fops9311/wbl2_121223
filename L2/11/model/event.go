package model

import (
	"errors"
	"sync"
	"time"
)

type Event struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (e Event) IsInInterval(start, end time.Time) bool {
	if start.Before(end) {
		return !e.Date.Before(start) && !e.Date.After(end)
	}
	if start.Equal(end) {
		return e.Date.Equal(start)
	}
	return false
}

type Calendar struct {
	events map[string]Event
	mu     sync.Mutex
}

func NewCalendar() *Calendar {
	return &Calendar{
		events: make(map[string]Event),
	}
}
func (c *Calendar) AddEvent(e Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.events[e.Id]; ok {
		return errors.New("already exists")
	}
	c.events[e.Id] = e
	return nil
}
func (c *Calendar) UpdateEvent(e Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.events[e.Id]; !ok {
		return errors.New("not exist")
	}
	c.events[e.Id] = e
	return nil
}
func (c *Calendar) DeleteEvent(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.events[id]; !ok {
		return errors.New("not exist")
	}
	delete(c.events, id)
	return nil
}
func (c *Calendar) GetEventsForDay(t time.Time) ([]Event, error) {
	var year, month, day = t.Year(), t.Month(), t.Day()
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour * 24)
	return c.GetTimeIntervalEvents(start, end), nil
}
func (c *Calendar) GetEventsForWeek(t time.Time) ([]Event, error) {
	var year, month, day = t.Year(), t.Month(), t.Day()
	var start, end time.Time
	switch t.Weekday() {
	case time.Monday:
		start = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.Month(month), day+7, 0, 0, 0, 0, time.UTC)
	case time.Tuesday:
		start = time.Date(year, time.Month(month), day-1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.Month(month), day+6, 0, 0, 0, 0, time.UTC)
	case time.Wednesday:
		start = time.Date(year, time.Month(month), day-2, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.Month(month), day+5, 0, 0, 0, 0, time.UTC)
	case time.Thursday:
		start = time.Date(year, time.Month(month), day-3, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.Month(month), day+4, 0, 0, 0, 0, time.UTC)
	case time.Friday:
		start = time.Date(year, time.Month(month), day-4, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.Month(month), day+3, 0, 0, 0, 0, time.UTC)
	case time.Saturday:
		start = time.Date(year, time.Month(month), day-5, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.Month(month), day+2, 0, 0, 0, 0, time.UTC)
	case time.Sunday:
		start = time.Date(year, time.Month(month), day-6, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.UTC)
	}
	return c.GetTimeIntervalEvents(start, end), nil
}
func (c *Calendar) GetEventsForMonth(t time.Time) ([]Event, error) {
	var year, month = t.Year(), t.Month()
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	var endyear, endmonth = nextMonth(year, int(month))
	end := time.Date(endyear, time.Month(endmonth), 1, 0, 0, 0, 0, time.UTC)
	return c.GetTimeIntervalEvents(start, end), nil
}
func nextMonth(year, month int) (int, int) {
	switch time.Month(month) {
	case time.December:
		return year + 1, month + 1
	default:
		return year, month + 1
	}
}
func (c *Calendar) GetEventsForYear(t time.Time) ([]Event, error) {
	var year = t.Year()
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Nanosecond * -1)
	return c.GetTimeIntervalEvents(start, end), nil
}
func (c *Calendar) GetTimeIntervalEvents(start, end time.Time) []Event {
	result := make([]Event, 0)
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, e := range c.events {
		if e.IsInInterval(start, end) {
			result = append(result, e)
		}
	}
	return result
}
