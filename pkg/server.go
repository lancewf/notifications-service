package pkg

import (
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Server
type Server struct {
	port int
}

func New(port int) Server {
	return Server{
		port: port,
	}
}

func (server Server) Start() {
	mutex := sync.Mutex{}

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		if r.Method == "POST" {
			log.Info("run POST")
			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else if r.Method == "GET" {

			log.Info("run GET")

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else {
			log.Info("Unhandled")
		}

	})
	http.HandleFunc("/inspec", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		if r.Method == "POST" {
			log.Info("inspec POST")
			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else if r.Method == "GET" {

			log.Info("inspec GET")

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte("{}"))
		} else {
			log.Info("Unhandled")
		}

	})
	http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)
}
