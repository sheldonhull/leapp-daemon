package controller

import (
  "github.com/gin-gonic/gin"
  logging2 "leapp_daemon/infrastructure/logging"
  websocket2 "leapp_daemon/infrastructure/websocket"
  "log"
  "net/http"
)

func WsController(context *gin.Context) {
	logging2.SetContext(context)
	serveWs(context.Writer, context.Request)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket2.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &websocket2.Connection{Send: make(chan websocket2.Message, 1), Ws: ws}
	s := websocket2.Subscription{Conn: c, Room: websocket2.DefaultRoom}

	websocket2.Hub.Register <- s

	go s.WritePump()
	go s.ReadPump()
}
