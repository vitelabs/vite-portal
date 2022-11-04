package handler

import "sync"

type Pubsub struct {
	topic string
	impl  *PubsubTopics
}

func NewPubsub() *Pubsub {
	ps := &Pubsub{
		topic: "",
	}
	ps.impl = NewPubsubTopics()
	return ps
}

func (ps *Pubsub) Subscribe() <-chan string {
	return ps.impl.Subscribe(ps.topic)
}

func (ps *Pubsub) Publish(msg string) {
	ps.impl.Publish(ps.topic, msg)
}

func (ps *Pubsub) Close() {
	ps.impl.Close()
}

type PubsubTopics struct {
	mu     sync.RWMutex
	subs   map[string][]chan string
	closed bool
}

func NewPubsubTopics() *PubsubTopics {
	ps := &PubsubTopics{}
	ps.subs = make(map[string][]chan string)
	return ps
}

func (ps *PubsubTopics) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *PubsubTopics) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

func (ps *PubsubTopics) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}
