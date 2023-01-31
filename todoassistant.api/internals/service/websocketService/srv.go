package websocketService

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketSrv interface {
	Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) (*websocket.Conn, error)
}

type websocketSrv struct{}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}


func (ws *websocketSrv) Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewWebSocketSrv() WebSocketSrv {
	return &websocketSrv{}
}
