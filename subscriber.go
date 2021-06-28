package pubsub

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Subscribers map[string]*Subscriber

type SubTopics map[string]bool

type Subscriber struct {
	id            string
	messages      chan *Message
	creation_time creationTime
	destroyed     bool
	topics        SubTopics
	lock          sync.RWMutex
}

func NewSubscriber() (*Subscriber, error) {
	id := make([]byte, 64)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}

	return &Subscriber{
		id:            hex.EncodeToString(id),
		messages:      make(chan *Message),
		creation_time: creationTime(time.Now().UnixNano()),
		destroyed:     false,
		lock:          sync.RWMutex{},
		topics:        SubTopics{},
	}, nil
}

func (s *Subscriber) GetID() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.id
}

func (s *Subscriber) GetCreationTimeAt() creationTime {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.creation_time
}

func (s *Subscriber) AddTopic(topic string) {
	s.lock.Lock()
	s.topics[topic] = true
	s.lock.Unlock()
}

func (s *Subscriber) DeleteTopic(topic string) {
	s.lock.Lock()
	delete(s.topics, topic)
	s.lock.Unlock()
}

func (s *Subscriber) GetTopics() []string {
	s.lock.RLock()
	subscriberTopics := s.topics
	s.lock.RUnlock()

	topics := []string{}
	for topic := range subscriberTopics {
		topics = append(topics, topic)
	}
	return topics
}

func (s *Subscriber) GetMessages() <-chan *Message {
	return s.messages
}

func (s *Subscriber) Deque(m *Message) *Subscriber {
	s.lock.RLock()
	if !s.destroyed {
		s.messages <- m
	}
	s.lock.RUnlock()

	return s
}

func (s *Subscriber) destroy() {
	s.lock.Lock()
	s.destroyed = true
	s.lock.Unlock()

	close(s.messages)
}

func (s *Subscriber) Ack(id string, ch <-chan *Message) {
	fmt.Printf("Subscriber_id %v, receiving...\n", id)
	for {
		if msg, ok := <-ch; ok {
			fmt.Printf("Subscriber_id %v, received: %v\n", id, msg.GetData())
		}
	}
}
