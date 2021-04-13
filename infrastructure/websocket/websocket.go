package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a Message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong Message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum Message size allowed from peer.
	maxMessageSize = 512
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Hub = hub{
	Rooms:      make(map[string]map[*Connection]bool),
	Register:   make(chan Subscription),
	Unregister: make(chan Subscription),
}

var DefaultRoom = "1"

// MessageType enum
type MessageType int
const (
	MfaTokenRequest MessageType = iota
)

type hub struct {
	Rooms      map[string]map[*Connection]bool
	Register   chan Subscription
	Unregister chan Subscription
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.Register:
			connections := h.Rooms[s.Room]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Rooms[s.Room] = connections
			}
			h.Rooms[s.Room][s.Conn] = true
			log.Println("registered", s)
			log.Println("updated h:", Hub)
		case s := <-h.Unregister:
			connections := h.Rooms[DefaultRoom]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.Send)
					if len(connections) == 0 {
						delete(h.Rooms, s.Room)
					}
				}
			}
			log.Println("unregistered", s)
			log.Println("updated h:", Hub)
		}
	}
}

type Connection struct {
	Ws   *websocket.Conn
	Send chan Message
}

func (c *Connection) write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}

type Message struct {
	MessageType MessageType
	Data        string
}

type Subscription struct {
	Conn *Connection
	Room string
}

func (s Subscription) ReadPump() {
	c := s.Conn

	defer func() {
		Hub.Unregister <- s
		c.Ws.Close()
	}()

	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Ws.SetPongHandler(func(string) error {
		log.Println("pong")
		log.Println("updated h:", Hub)
		c.Ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.Ws.SetCloseHandler(func(int, string) error {
		c.Ws.Close()
		return nil
	})

	for {
		log.Println("waiting for a Err...")
		_, msgByteArray, err := c.Ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		//m := Err{msg, []byte(DefaultRoom)}
		var msg Message
		err = json.Unmarshal(msgByteArray, &msg)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		switch msg.MessageType {
		case MfaTokenRequest:

		default:
			// do nothing
		}
	}
}

func (s *Subscription) WritePump() {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			messageJson, err := json.Marshal(message)
			if err != nil {
				return
			}
			if err := c.write(websocket.TextMessage, messageJson); err != nil {
				return
			}
		case <-ticker.C:
			log.Println("ping")
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func SendMessage(msg Message) error {
	connections := Hub.Rooms[DefaultRoom]

	for c := range connections {
		select {
		case c.Send <- msg:
		default:
			close(c.Send)
			delete(connections, c)
			if len(connections) == 0 {
				delete(Hub.Rooms, DefaultRoom)
			}
		}
	}

	return nil
}

type MfaTokenRequestData struct {
	SessionId string
}
