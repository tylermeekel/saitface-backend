package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/olahol/melody"
)

type Server struct {
	DB          *sql.DB
	Melody      *melody.Melody
	RestyClient *resty.Client
}

func (s *Server) RunServer() {
	mux := chi.NewMux()

	s.Melody = s.NewMelody()

	rc := resty.New()

	s.RestyClient = rc

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	s.DB = db

	s.QueryDeleteOldThreads()
	go s.DeleteOldThreadsTick()

	mux.Get("/ws", s.WrapMelody)
	mux.Mount("/threads", s.ThreadsRouter())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("Listening on port", port)
	http.ListenAndServe(":"+port, mux)
}

func (s *Server) WrapMelody(w http.ResponseWriter, r *http.Request) {
	err := s.Melody.HandleRequest(w, r)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) DeleteOldThreadsTick() {
	for range time.Tick(5 * time.Minute) {
		s.QueryDeleteOldThreads()
	}
}

func (s *Server) QueryDeleteOldThreads() {
	fmt.Println("Running Delete Old Threads Command")
	queryRows, err := s.DB.Query("DELETE FROM threads WHERE last_bumped < NOW() - interval '5 minute' RETURNING id")
	if err != nil {
		fmt.Println(err)
	}

	var ids []int

	for queryRows.Next() {
		var id int
		queryRows.Scan(&id)
		ids = append(ids, id)
	}

	fmt.Printf("Removed %d rows\n", len(ids))
}
