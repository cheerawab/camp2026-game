package matches

import (
	"sync"

	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

type Event struct {
	Name    string
	Match   mongomodel.Match
	Answers []mongomodel.MatchAnswer
}

type Broker struct {
	mu          sync.Mutex
	subscribers map[string]map[chan Event]struct{}
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string]map[chan Event]struct{}),
	}
}

func (b *Broker) Subscribe(matchID string) (<-chan Event, func()) {
	ch := make(chan Event, 8)

	b.mu.Lock()
	if _, ok := b.subscribers[matchID]; !ok {
		b.subscribers[matchID] = make(map[chan Event]struct{})
	}
	b.subscribers[matchID][ch] = struct{}{}
	b.mu.Unlock()

	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		if subscribers, ok := b.subscribers[matchID]; ok {
			delete(subscribers, ch)
			if len(subscribers) == 0 {
				delete(b.subscribers, matchID)
			}
		}
		close(ch)
	}

	return ch, unsubscribe
}

func (b *Broker) Publish(matchID string, event Event) {
	b.mu.Lock()
	subscribers := make([]chan Event, 0, len(b.subscribers[matchID]))
	for ch := range b.subscribers[matchID] {
		subscribers = append(subscribers, ch)
	}
	b.mu.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
