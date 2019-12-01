package workers

import (
	"github.com/gorilla/mux"
	"github.com/pojntfx/godibs/pkg/utils"
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

// Start starts a GitHTTPWorker
func (worker *GitHTTPWorker) Start(errors chan error, events chan utils.Event) {
	events <- utils.Event{
		Code:    0,
		Message: "Started",
	}

	r := mux.NewRouter()

	r.PathPrefix(worker.HTTPPathPrefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		events <- utils.Event{
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

	events <- utils.Event{
		Code:    2,
		Message: "Server stopped",
	}
}