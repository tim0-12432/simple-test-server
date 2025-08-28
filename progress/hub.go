package progress

import (
	"sync"
	"time"
)

type Event struct {
	Percent int    `json:"percent"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type Hub struct {
	mu    sync.Mutex
	chans map[string]chan Event
}

var Default = &Hub{chans: map[string]chan Event{}}

func (h *Hub) New(id string) chan Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan Event, 16)
	if old, ok := h.chans[id]; ok {
		close(old)
	}
	h.chans[id] = ch
	go func() {
		time.Sleep(10 * time.Minute)
		h.Remove(id)
	}()
	return ch
}

func (h *Hub) Get(id string) (chan Event, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch, ok := h.chans[id]
	return ch, ok
}

func (h *Hub) Remove(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.chans[id]; ok {
		close(ch)
		delete(h.chans, id)
	}
}

func (h *Hub) Send(id string, e Event) {
	h.mu.Lock()
	ch, ok := h.chans[id]
	h.mu.Unlock()
	if !ok {
		return
	}
	select {
	case ch <- e:
	default:
		// drop if full
	}
}
