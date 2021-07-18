package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

var tests = []struct {
	name                string
	in                  storage.Event
	function            string
	expectedError       error
	expectedNumOfEvents int
}{
	{
		name: "simple test",
		in: storage.Event{
			ID:    "first event",
			Title: "it is title",
		},
		function:      "CreateEvent",
		expectedError: nil,
	},
	{
		name: "test for duplication",
		in: storage.Event{
			ID: "first event",
		},
		function:      "CreateEvent",
		expectedError: ErrEventIDIsExist,
	},
	{
		name: "trying to update nonexistent event",
		in: storage.Event{
			ID: "not exist",
		},
		function:      "UpdateEvent",
		expectedError: ErrEventIDIsNotExist,
	},
	{
		name: "trying to update existent event",
		in: storage.Event{
			ID:    "first event",
			Title: "this title is changed",
		},
		function:      "UpdateEvent",
		expectedError: nil,
	},
	{
		name: "checking that event is changed",
		in: storage.Event{
			ID:    "first event",
			Title: "this title is changed",
		},
		function:      "getEvent",
		expectedError: nil,
	},
	{
		name: "checking getEvent function for nonexistent event",
		in: storage.Event{
			ID: "not exist",
		},
		function:      "getEvent",
		expectedError: ErrEventIDIsNotExist,
	},
	{
		name: "checking of DeleteEvent function. 1/3 deleting of nonexistent event",
		in: storage.Event{
			ID: "not exist",
		},
		function:      "DeleteEvent",
		expectedError: ErrEventIDIsNotExist,
	},
	{
		name: "checking of DeleteEvent function. 2/3 deleting of existent event",
		in: storage.Event{
			ID: "first event",
		},
		function:      "DeleteEvent",
		expectedError: nil,
	},
	{
		name: "checking of DeleteEvent function. 3/3 check that deletion is complete",
		in: storage.Event{
			ID: "first event",
		},
		function:      "getEvent",
		expectedError: ErrEventIDIsNotExist,
	},
	{
		name: "preparing for GetEventsFor... functions. 1/4 creating 1 of 2 events for a day, week and month",
		in: storage.Event{
			ID:       "1/2 event of day",
			DateTime: time.Date(2021, time.September, 15, 12, 0, 0, 0, time.Local),
		},
		function:      "CreateEvent",
		expectedError: nil,
	},
	{
		name: "preparing for GetEventsFor... functions. 2/4 creating 2 of 2 events for a day, week and month",
		in: storage.Event{
			ID:       "2/2 event of day",
			DateTime: time.Date(2021, time.September, 15, 13, 30, 0, 0, time.Local),
		},
		function:      "CreateEvent",
		expectedError: nil,
	},
	{
		name: "preparing for GetEventsFor... functions. 3/4 creating event for a week and month",
		in: storage.Event{
			ID:       "event of week",
			DateTime: time.Date(2021, time.September, 17, 23, 45, 0, 0, time.Local),
		},
		function:      "CreateEvent",
		expectedError: nil,
	},
	{
		name: "preparing for GetEventsFor... functions. 4/4 creating event for a month",
		in: storage.Event{
			ID:       "event of month",
			DateTime: time.Date(2021, time.September, 22, 23, 45, 0, 0, time.Local),
		},
		function:      "CreateEvent",
		expectedError: nil,
	},
	{
		name: "checking of GetEventsForDay function",
		in: storage.Event{
			DateTime: time.Date(2021, time.September, 15, 0, 0, 0, 0, time.Local),
		},
		function:            "GetEventsForDay",
		expectedError:       nil,
		expectedNumOfEvents: 2,
	},
	{
		name: "checking of GetEventsForWeek function",
		in: storage.Event{
			DateTime: time.Date(2021, time.September, 15, 0, 0, 0, 0, time.Local),
		},
		function:            "GetEventsForWeek",
		expectedError:       nil,
		expectedNumOfEvents: 3,
	},
	{
		name: "checking of GetEventsForMonth function",
		in: storage.Event{
			DateTime: time.Date(2021, time.September, 15, 0, 0, 0, 0, time.Local),
		},
		function:            "GetEventsForMonth",
		expectedError:       nil,
		expectedNumOfEvents: 4,
	},
}

func TestForBuncOfCases(t *testing.T) {
	s := New()
	for _, tt := range tests {
		t.Run(fmt.Sprintf("subtest: %v", tt.name), func(t *testing.T) {
			tt := tt
			var errFromFunc error
			switch tt.function {
			case "CreateEvent":
				errFromFunc = s.CreateEvent(context.TODO(), tt.in)
			case "UpdateEvent":
				errFromFunc = s.UpdateEvent(context.TODO(), tt.in)
			case "getEvent":
				var event storage.Event
				event, errFromFunc = s.getEvent(tt.in.ID)
				if errFromFunc == nil {
					require.Equal(t, tt.in, event)
				}
			case "DeleteEvent":
				errFromFunc = s.DeleteEvent(context.TODO(), tt.in)
			case "GetEventsForDay":
				var events []storage.Event
				events, errFromFunc = s.GetEventsForDay(context.TODO(), tt.in.DateTime)
				require.Equal(t, tt.expectedNumOfEvents, len(events))
			case "GetEventsForWeek":
				var events []storage.Event
				events, errFromFunc = s.GetEventsForWeek(context.TODO(), tt.in.DateTime)
				require.Equal(t, tt.expectedNumOfEvents, len(events))
			case "GetEventsForMonth":
				var events []storage.Event
				events, errFromFunc = s.GetEventsForMonth(context.TODO(), tt.in.DateTime)
				require.Equal(t, tt.expectedNumOfEvents, len(events))
			default:
				panic("you should add new cases or fill tt.function with right value")
			}
			switch {
			case errors.Is(tt.expectedError, ErrEventIDIsExist):
				require.ErrorIs(t, errFromFunc, ErrEventIDIsExist, "need: %v, got: %v", ErrEventIDIsExist, errFromFunc)
			case errors.Is(tt.expectedError, ErrEventIDIsNotExist):
				require.ErrorIs(t, errFromFunc, ErrEventIDIsNotExist, "need: %v, got: %v", ErrEventIDIsNotExist, errFromFunc)
			case tt.expectedError == nil:
				require.NoErrorf(t, errFromFunc, "need: no error, got: ", errFromFunc)
			default:
				panic("you should add this error in tt.expectedError")
			}
		})
	}
}

func TestForConcurrency(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("making many events by many workers instantly", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		storageMain := New()
		countOfWorkers := 50
		countOfEvents := 100
		wg.Add(countOfWorkers)
		for worker := 0; worker < countOfWorkers; worker++ {
			worker := worker
			go func() {
				defer wg.Done()
				for event := 0; event < countOfEvents; event++ {
					event := event
					err := storageMain.CreateEvent(context.TODO(), storage.Event{
						ID: storage.EventID(fmt.Sprintf("event #%v from worker #%v", event, worker)),
					})
					require.NoError(t, err, "expected no errors")
				}
			}()
		}
		wg.Wait()
		require.Equal(t, countOfWorkers*countOfEvents, len(storageMain.Events))

		wg.Add(countOfWorkers)
		for worker := 0; worker < countOfWorkers; worker++ {
			worker := worker
			go func() {
				defer wg.Done()
				for event := 0; event < countOfEvents; event++ {
					event := event
					err := storageMain.DeleteEvent(context.TODO(), storage.Event{
						ID: storage.EventID(fmt.Sprintf("event #%v from worker #%v", event, worker)),
					})
					require.NoError(t, err, "expected no errors")
				}
			}()
		}
		wg.Wait()
		require.Equal(t, 0, len(storageMain.Events))
	})
}
