package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"saitface/internal/utils"
	"strconv"
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
	r.Get("/{id}", s.GetOneThread)
	r.Post("/", s.CreateNewThread)

	return r
}

func (s *Server) QueryOneThread(id int) (Thread, error) {
	var thread Thread

	row := s.DB.QueryRow("SELECT * FROM threads WHERE id=$1", id)

	err := row.Scan(&thread.ID, &thread.Title, pq.Array(&thread.Interests), &thread.LastBumped)

	return thread, err
}

func (s *Server) GetOneThread(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Incorrect ID Value"))
		return
	}

	thread, err := s.QueryOneThread(id)
	if err != nil {
		fmt.Println(err)
	}

	utils.SendJSON(w, thread)
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

	var interestStruct struct {
		Interests []string `json:"interests"`
	}

	interestStruct.Interests = interests

	b, err := json.Marshal(interestStruct)
	if err != nil {
		return thread, errors.New("Error marshalling interests")
	}

	title := s.GetThreadTitle(string(b))

	row := s.DB.QueryRow("INSERT INTO threads(title, interests) VALUES($1, $2) RETURNING *", title, pq.Array(interests))

	err = row.Scan(&thread.ID, &thread.Title, pq.Array(&thread.Interests), &thread.LastBumped)

	return thread, err
}

func (s *Server) GetThreadTitle(query string) string {
	fmt.Println(query)
	var response struct {
		Title string `json:"title"`
	}

	tsurl := os.Getenv("TITLE_SERVER_URL")
	fmt.Println(tsurl)

	resp, err := s.RestyClient.R().
		SetBody(query).
		Post(tsurl)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	str := resp.String()
	fmt.Println(str)

	json.Unmarshal([]byte(str), &response)

	fmt.Println(response.Title)
	return response.Title
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

func (s *Server) QueryBumpThread(id int) {
	currentTime := time.Now()

	row := s.DB.QueryRow("UPDATE threads SET last_bumped=$1 WHERE id=$2 RETURNING id", currentTime, id)

	var returnedID int
	err := row.Scan(&returnedID)
	if err != nil {
		fmt.Println("Error bumping thread", err)
	}
}
