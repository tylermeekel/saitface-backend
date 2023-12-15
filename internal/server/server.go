package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/olahol/melody"
)

type Server struct {
	DB *sql.DB
	Melody *melody.Melody
}

func (s *Server) RunServer() { 
	mux := chi.NewMux()

	s.Melody = NewMelody()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil{
		log.Fatalln(err)
	}

	s.DB = db

	mux.Get("/ws", s.WrapMelody)
	mux.Mount("/threads", s.ThreadsRouter())

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", mux)
}

func (s *Server) WrapMelody(w http.ResponseWriter, r *http.Request) {
	err := s.Melody.HandleRequest(w, r)
	if err != nil{
		fmt.Println(err)
	}
}