package app

import (
	"context"
	"fmt"
	"time"

	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  logger.Logger
	Storage storage.Storage
}

// New returns a new App main object.
func New(logger logger.Logger, storage storage.Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

// CreateEvent adds event into storage.
func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	fmt.Printf("%#v", a.Storage)
	return a.Storage.CreateEvent(ctx, event)
}

// UpdateEvent updates event in storage.
func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.Storage.UpdateEvent(ctx, event)
}

// DeleteEvent deletes event from storage.
func (a *App) DeleteEvent(ctx context.Context, event storage.Event) error {
	return a.Storage.DeleteEvent(ctx, event)
}

// GetEventsForDay return slice of events for one day.
func (a *App) GetEventsForDay(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	return a.Storage.GetEventsForDay(ctx, dateOfQuestion)
}

// GetEventsForWeek returns slice of events for one week, you need to determine any day from that week.
func (a *App) GetEventsForWeek(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	return a.Storage.GetEventsForWeek(ctx, dateOfQuestion)
}

// GetEventsForMonth returns slice of events for one month, you need to determine any day from that month.
func (a *App) GetEventsForMonth(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	return a.Storage.GetEventsForMonth(ctx, dateOfQuestion)
}

// Debug make a message in logger at Debug-level.
func (a *App) Debug(msg string) {
	a.Logger.Debug(msg)
}

// Info make a message in logger at Info-level.
func (a *App) Info(msg string) {
	a.Logger.Info(msg)
}

// Warn make a message in logger at Warn-level.
func (a *App) Warn(msg string) {
	a.Logger.Warn(msg)
}

// Error make a message in logger at Error-level.
func (a *App) Error(msg string) {
	a.Logger.Error(msg)
}
