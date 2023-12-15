package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Thread struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Interests []string `json:"interests"`
}

// Available at [url]/threads
func (s *Server) ThreadsRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", s.GetAllThreads)

	return r
}

func (s *Server) QueryAllThreads() ([]Thread, error) {
	var threads []Thread
	rows, err := s.DB.Query("SELECT * FROM threads")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var thread Thread

		rows.Scan(&thread)
		threads = append(threads, thread)
	}

	return threads, nil
}

func (s *Server) GetAllThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := s.QueryAllThreads()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	json.NewEncoder(w).Encode(threads)
}
