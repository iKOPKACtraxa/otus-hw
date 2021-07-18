package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventIDIsExist    = errors.New("ivent ID is already exists")
	ErrEventIDIsNotExist = errors.New("ivent is not exists")
)

type Storageinmemory struct {
	mu     sync.RWMutex
	Events map[storage.EventID]storage.Event
}

func New() *Storageinmemory {
	return &Storageinmemory{
		Events: make(map[storage.EventID]storage.Event),
	}
}

// CreateEvent adds event into storage.
func (s *Storageinmemory) CreateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Events[event.ID]; ok {
		return ErrEventIDIsExist
	}
	s.Events[event.ID] = event
	return nil
}

// UpdateEvent updates event in storage.
func (s *Storageinmemory) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Events[event.ID]; !ok {
		return ErrEventIDIsNotExist
	}
	s.Events[event.ID] = event
	return nil
}

// DeleteEvent deletes event from storage.
func (s *Storageinmemory) DeleteEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Events[event.ID]; !ok {
		return ErrEventIDIsNotExist
	}
	delete(s.Events, event.ID)
	return nil
}

// GetEventsForDay return slice of events for one day.
func (s *Storageinmemory) GetEventsForDay(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var sliceOfEvents []storage.Event
	for _, v := range s.Events {
		if v.DateTime.Year() == dateOfQuestion.Year() && v.DateTime.YearDay() == dateOfQuestion.YearDay() {
			sliceOfEvents = append(sliceOfEvents, v)
		}
	}
	return sliceOfEvents, nil
}

// GetEventsForWeek returns slice of events for one week, you need to determine any day from that week.
func (s *Storageinmemory) GetEventsForWeek(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	year, week := dateOfQuestion.ISOWeek()
	var sliceOfEvents []storage.Event
	for _, v := range s.Events {
		if y, w := v.DateTime.ISOWeek(); y == year && w == week {
			sliceOfEvents = append(sliceOfEvents, v)
		}
	}
	return sliceOfEvents, nil
}

// GetEventsForMonth returns slice of events for one month, you need to determine any day from that month.
func (s *Storageinmemory) GetEventsForMonth(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	year, month, _ := dateOfQuestion.Date()
	var sliceOfEvents []storage.Event
	for _, v := range s.Events {
		if y, m, _ := v.DateTime.Date(); y == year && m == month {
			sliceOfEvents = append(sliceOfEvents, v)
		}
	}
	return sliceOfEvents, nil
}

// getEvent is for test purpose, by ID of event it returns event from storage.
func (s *Storageinmemory) getEvent(eventID storage.EventID) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Events[eventID]; !ok {
		return storage.Event{}, ErrEventIDIsNotExist
	}
	return s.Events[eventID], nil
}
