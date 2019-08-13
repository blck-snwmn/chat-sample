package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//originのチェックは行わない
		return true
	},
}

func read(room *Room, conn *websocket.Conn) {
	defer conn.Close()
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		mc := MessageContainer{
			msgType: msgType,
			message: msg,
		}
		room.ch <- mc
	}
}

func write(room *Room, conn *websocket.Conn) {
	defer conn.Close()
	for {
		select {
		case mc := <-room.ch:
			for _, cn := range room.conns {
				if err := cn.WriteMessage(mc.msgType, mc.message); err != nil {
					log.Println(err)
					return
				}
			}
		default:
		}
	}
}

func webSocketHandler(room *Room, w http.ResponseWriter, r *http.Request) {
	log.Println("in")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	room.conns = append(room.conns, conn)
	go read(room, conn)
	go write(room, conn)
}
func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// MessageContainer is container of message webscoket recieved
type MessageContainer struct {
	msgType int
	message []byte
}

//Room is chat room
type Room struct {
	conns []*websocket.Conn
	ch    chan MessageContainer
}

func main() {
	room := Room{
		conns: []*websocket.Conn{},
		ch:    make(chan MessageContainer),
	}
	var httpServer http.Server
	httpServer.Addr = ":28888"
	http.HandleFunc("/ws", handler)
	http.HandleFunc("/ws/send", func(w http.ResponseWriter, r *http.Request) {
		webSocketHandler(&room, w, r)
	})
	httpServer.ListenAndServe()
}
