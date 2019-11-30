package workers

import (
	"github.com/gorilla/mux"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"gopkg.in/mysticmode/gitviahttp.v1"
	"net/http"
	"strconv"
	"sync"
)

// GitHTTPWorker serves Git repos via HTTP
type GitHTTPWorker struct {
	ReposDir       string // Directory in which the repos that should be served reside
	HTTPPathPrefix string // Path prefix for the HTTP server
	Port           int    // Port on which the HTTP server should listen
}

// Start starts a GitHTTPWorker
func (config *GitHTTPWorker) Start(wg *sync.WaitGroup) {
	log.Info("Starting Git HTTP worker ...", rz.String("ReposDir", config.ReposDir), rz.String("HTTPPathPrefix", config.HTTPPathPrefix), rz.Int("Port", config.Port))

	r := mux.NewRouter()

	r.PathPrefix(config.HTTPPathPrefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gitviahttp.Context(w, r, config.ReposDir)
	}).Methods("GET", "POST")

	s := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + strconv.Itoa(config.Port),
	}

	s.ListenAndServe()

	wg.Done()
}
