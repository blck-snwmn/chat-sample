package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ChatRoot is Chat's root
type ChatRoot struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func (root *ChatRoot) getRoom(name string) *Room {
	root.mu.Lock()
	defer root.mu.Unlock()

	room, exist := root.rooms[name]
	if !exist {
		room = &Room{
			conns:    make([]*websocket.Conn, 0),
			recieved: make(chan MessageContainer),
		}
		root.rooms[name] = room
	}
	return room
}

func newChatRoom() *ChatRoot {
	return &ChatRoot{rooms: make(map[string]*Room)}
}
