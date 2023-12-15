package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"saitface/internal/utils"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

type Thread struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Interests  []string  `json:"interests"`
	LastBumped time.Time `json:"last_bumped,omitempty"`
}

// Available at [url]/threads
func (s *Server) ThreadsRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", s.GetAllThreads)
	r.Post("/", s.CreateNewThread)

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

		rows.Scan(&thread.ID, &thread.Title, pq.Array(&thread.Interests), &thread.LastBumped)
		threads = append(threads, thread)
	}

	return threads, nil
}

func (s *Server) QueryNewThread(interests []string) (Thread, error) {
	var thread Thread

	if len(interests) < 1 {
		return thread, errors.New("Error: Not enough interests, need at least 1")
	}

	query := strings.Join(interests, "-")

	title := GetThreadTitle(query)

	row := s.DB.QueryRow("INSERT INTO threads(title, interests) VALUES($1, $2) RETURNING *", title, pq.Array(interests))

	err := row.Scan(&thread.ID, &thread.Title, pq.Array(&thread.Interests), &thread.LastBumped)

	return thread, err
}

func GetThreadTitle(query string) string {
	return "Title"
}

func (s *Server) CreateNewThread(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Interests []string `json:"interests"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	thread, err := s.QueryNewThread(req.Interests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	utils.SendJSON(w, thread)
}

func (s *Server) GetAllThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := s.QueryAllThreads()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	utils.SendJSON(w, threads)
}
