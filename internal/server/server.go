package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	s.Melody = NewMelody()

	rc := resty.New()

	s.RestyClient = rc

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	s.DB = db

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
