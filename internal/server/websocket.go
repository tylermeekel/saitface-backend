package server

import (
	"encoding/json"
	"fmt"

	"github.com/olahol/melody"
)

type WSMessageType int

const (
	Connect WSMessageType = iota
	SendMessage
)

type WebSocketMessage struct {
	Type     WSMessageType `json:"type"`
	ThreadID int           `json:"thread_id"`
	Username string        `json:"username"`
	Content  string        `json:"content"`
}

func (server *Server) NewMelody() *melody.Melody {
	m := melody.New()

	m.HandleConnect(func(s *melody.Session) {
		s.Write([]byte("connected"))
	})

	m.HandleMessage(func(s *melody.Session, data []byte) {
		var message WebSocketMessage

		json.Unmarshal(data, &message)

		switch message.Type {
		case Connect:
			s.Set("thread", message.ThreadID)
			fmt.Println("Set the ID")
		case SendMessage:
			server.QueryBumpThread(message.ThreadID)
			m.BroadcastFilter(data, func(checkedSession *melody.Session) bool {
				threadVal, exists := checkedSession.Get("thread")
				if !exists {
					return false
				}

				checkedThreadID := threadVal.(int)

				if checkedThreadID == message.ThreadID && checkedSession != s {
					return true
				} else {
					return false
				}
			})
		}
	})

	return m
}
