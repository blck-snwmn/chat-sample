package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestReturnSameRoom(t *testing.T) {
	root := newChatRoom()
	room := root.getRoom("test")
	if room != root.getRoom("test") {
		t.Error("expected same room")
	}
}

func TestReturnNewRoom(t *testing.T) {
	root := newChatRoom()
	room := root.getRoom("test")
	if room == root.getRoom("other") {
		t.Error("expected different room")
	}
}

func TestMutex(t *testing.T) {
	expectedLen := 100

	root := newChatRoom()

	wg := &sync.WaitGroup{}
	for i := 0; i < expectedLen; i++ {
		wg.Add(1)
		go func(counter int) {
			defer wg.Done()
			root.getRoom("test" + strconv.Itoa(counter))
		}(i)
	}
	wg.Wait()
	if len(root.rooms) != expectedLen {
		t.Errorf("rooms length is %d; want %d.", len(root.rooms), expectedLen)
	}
}
