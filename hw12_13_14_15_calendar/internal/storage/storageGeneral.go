package storage

import (
	"context"
	"time"
)

type (
	EventID string
	Event   struct {
		ID          EventID       `json:"ID"`          // ID - уникальный идентификатор события (можно воспользоваться UUID);
		Title       string        `json:"Title"`       // Заголовок - короткий текст;
		DateTime    time.Time     `json:"DateTime"`    // Дата и время события;
		Duration    time.Duration `json:"Duration"`    // Длительность события (или дата и время окончания);
		Description string        `json:"Description"` // Описание события - длинный текст, опционально;
		UserID      string        `json:"UserID"`      // ID пользователя, владельца события;
		NoteBefore  time.Duration `json:"NoteBefore"`  // За сколько времени высылать уведомление, опционально.
	}
)

type Storage interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, event Event) error
	GetEventsForDay(ctx context.Context, dateOfQuestion time.Time) ([]Event, error)
	GetEventsForWeek(ctx context.Context, dateOfQuestion time.Time) ([]Event, error)
	GetEventsForMonth(ctx context.Context, dateOfQuestion time.Time) ([]Event, error)
}
