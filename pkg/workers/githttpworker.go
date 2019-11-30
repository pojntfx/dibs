package workers

import (
	"github.com/gorilla/mux"
	"gopkg.in/mysticmode/gitviahttp.v1"
	"net/http"
	"strconv"
)

// GitHTTPWorker serves Git repos via HTTP
type GitHTTPWorker struct {
	ReposDir       string // Directory in which the repos that should be served reside
	HTTPPathPrefix string // Path prefix for the HTTP server
	Port           int    // Port on which the HTTP server should listen
}

// GitHTTPWorkerEvent enables status messages
type GitHTTPWorkerEvent struct {
	Code    int    // Status code of the event
	Message string // Message of the event
}

// Start starts a GitHTTPWorker
func (worker *GitHTTPWorker) Start(errors chan error, events chan GitHTTPWorkerEvent) {
	events <- GitHTTPWorkerEvent{
		Code:    0,
		Message: "Started",
	}

	r := mux.NewRouter()

	r.PathPrefix(worker.HTTPPathPrefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		events <- GitHTTPWorkerEvent{
			Code:    1,
			Message: r.Method + " request to " + r.URL.Path + " received",
		}

		gitviahttp.Context(w, r, worker.ReposDir)
	}).Methods("GET", "POST")

	s := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + strconv.Itoa(worker.Port),
	}

	if err := s.ListenAndServe(); err != nil {
		errors <- err
		return
	}

	events <- GitHTTPWorkerEvent{
		Code:    2,
		Message: "Server stopped",
	}
}
