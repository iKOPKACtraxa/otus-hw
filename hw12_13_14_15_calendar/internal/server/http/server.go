package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	HTTPServer *http.Server
	App        Application
}

type Application interface {
	storage.Storage
	logger.Logger
}

// ServeHTTP performs a requests and writes a response.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Это выводится в Terminal/Browser при выполнении ServeHTTP\n"))

	/* черновик, удалю
	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("Это выводится в Terminal/Browser при выполнении ServeHTTP\n"))
	var dm storage.Event                                                           // для теста инпута
	var body []byte                                                                // для теста инпута

	body, _ = ioutil.ReadAll(r.Body) // для теста инпута
	// for _, v := range body {
	// 	fmt.Print(string(v))
	// }
	err := json.Unmarshal(body, &dm) // для теста инпута
	if err != nil {
		s.App.Error("at unmarshalling has got an error: " + err.Error())
	}
	// fmt.Println(dm) // для теста инпута
	// s.App.CreateEvent(r.Context(), dm) // для теста инпута1
	// s.App.UpdateEvent(r.Context(), dm) 						 // для теста инпута2
	// s.App.GetEventsForDay(r.Context(), dm.DateTime) // для теста инпута3
	slice, err := s.App.GetEventsForMonth(r.Context(), dm.DateTime) // для теста инпута3

	// после запуска make run запустить следующую команду для отправки запроса в базу
	// 	curl --header "Content-Type: application/json"   --request GET   --data '{
	//   "DateTime": "2021-07-13T18:43:10+03:00"
	// }'   http://localhost:8081/
	if err != nil {
		s.App.Error("at request has got an error: " + err.Error())
	} else {
		for _, v := range slice {
			w.Write([]byte(fmt.Sprint(v)))
			w.Write([]byte("\n"))
		}
	}

	// w.Header().Set("Content-Type", "application/json") //todo сделать потом вывод json
	*/
}

// NewServer returns a new Server object.
func NewServer(app Application, hostPort string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		HTTPServer: &http.Server{
			Addr:         hostPort,
			Handler:      mux,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		App: app,
	}

	mux.Handle("/", LoggingMiddleware(app, server))
	return server
}

// Start starts Server.
func (s *Server) Start(ctx context.Context) error {
	err := s.HTTPServer.ListenAndServe()
	<-ctx.Done()
	return err
}

// Stop stops Server gracefully.
func (s *Server) Stop(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
