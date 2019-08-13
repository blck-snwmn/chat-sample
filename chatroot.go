package main

import "github.com/gorilla/websocket"

// ChatRoot is Chat's root
type ChatRoot struct {
	rooms map[string]*Room
}

func (root *ChatRoot) getRoom(name string) *Room {
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
