package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"leapp_daemon/logging"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func WsController(context *gin.Context) {
	logging.SetContext(context)

	roomId := context.Param("roomId")
	serveWs(context.Writer, context.Request, roomId)
}

func serveWs(w http.ResponseWriter, r *http.Request, roomId string) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}

	Hub.register <- s

	go s.writePump()
	go s.readPump()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type subscription struct {
	conn *connection
	room string
}

func (s subscription) readPump() {
	c := s.conn

	defer func() {
		Hub.unregister <- s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		log.Println("pong")
		log.Println("updated h:", Hub)
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.ws.SetCloseHandler(func(int, string) error {
		c.ws.Close()
		return nil
	})

	for {
		log.Println("waiting for a message...")
		_, msg, err := c.ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Println("message:", string(msg))

		m := message{msg, s.room}
		Hub.broadcast <- m
	}
}

func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
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

type connection struct {
	ws *websocket.Conn
	send chan []byte
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

type message struct {
	data []byte
	room string
}

type hub struct {
	rooms map[string]map[*connection]bool
	broadcast chan message
	register chan subscription
	unregister chan subscription
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
			log.Println("registered", s)
			log.Println("updated h:", Hub)
		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
			log.Println("unregistered", s)
			log.Println("updated h:", Hub)
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
			log.Println("broadcasted message", string(m.data))
			log.Println("updated h:", Hub)
		}
	}
}

var Hub = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}