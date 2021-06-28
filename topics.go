package pubsub

import (
	"fmt"
	"sync"
	"time"
)

type TopicsTable map[string]Subscribers

type Topics struct {
	reg_subscribers Subscribers
	sub_lock        sync.RWMutex
	sub_topics      TopicsTable
	topic_lock      sync.RWMutex
}

//create new topic with given ID
func NewTopic(topic_id string) *Topics {
	new_topic := Topics{
		reg_subscribers: Subscribers{},
		sub_lock:        sync.RWMutex{},
		sub_topics:      TopicsTable{},
		topic_lock:      sync.RWMutex{},
	}
	new_topic.CreateTopic(topic_id)
	return &new_topic
}

//Register and create new subscriber
func (b *Topics) Register_Sub() (*Subscriber, error) {
	s, err := NewSubscriber()

	if err != nil {
		return nil, err
	}

	b.sub_lock.Lock()
	b.reg_subscribers[s.GetID()] = s
	b.sub_lock.Unlock()

	return s, nil
}

//Create new Topic to the existing topic
func (b *Topics) CreateTopic(topic_id string) {
	b.sub_lock.Lock()
	if nil == b.sub_topics[topic_id] {
		b.sub_topics[topic_id] = Subscribers{}
	}
	b.sub_lock.Unlock()
}

//Delete existing topic
func (b *Topics) DeleteTopic(topic_id string) {
	b.sub_lock.Lock()
	for _, s := range b.sub_topics[topic_id] {
		s.DeleteTopic(topic_id)
	}
	if nil != b.sub_topics[topic_id] {
		delete(b.sub_topics, topic_id)
	}
	b.sub_lock.Unlock()
}

//Subscribe the given sub to topic list
func (b *Topics) Subscribe(s *Subscriber, topics ...string) {
	b.topic_lock.Lock()
	defer b.topic_lock.Unlock()
	for _, topic := range topics {
		if nil == b.sub_topics[topic] {
			fmt.Println("!!! topic does not exist !!!")
		}
		s.AddTopic(topic)
		b.sub_topics[topic][s.id] = s
	}
}

//Unsubscribe given sub to topic list
func (b *Topics) Unsubscribe(s *Subscriber, topics ...string) {
	for _, topic := range topics {
		b.topic_lock.Lock()
		if nil == b.sub_topics[topic] {
			b.topic_lock.Unlock()
			continue
		}
		delete(b.sub_topics[topic], s.id)
		b.topic_lock.Unlock()
		s.DeleteTopic(topic)
	}
}

//Completely remove the given subscriber
func (b *Topics) RemoveSubscriber(s *Subscriber) {
	s.destroy()
	b.sub_lock.Lock()
	b.Unsubscribe(s, s.GetTopics()...)
	delete(b.reg_subscribers, s.id)
	defer b.sub_lock.Unlock()
}

//Send data to all subs
func (b *Topics) Broadcast(data interface{}, topics ...string) {
	for _, topic := range topics {
		if b.Subscribers(topic) < 1 {
			continue
		}
		b.topic_lock.RLock()
		for _, s := range b.sub_topics[topic] {
			m := &Message{
				topic_name:    Topic(topic),
				data:          data,
				creation_time: creationTime(time.Now().UnixNano()),
			}
			go (func(s *Subscriber) {
				s.Deque(m)
			})(s)
		}
		b.topic_lock.RUnlock()
	}
}

func (b *Topics) Subscribers(topic string) int {
	b.topic_lock.RLock()
	defer b.topic_lock.RUnlock()
	return len(b.sub_topics[topic])
}

func (b *Topics) GetTopics() []string {
	b.topic_lock.RLock()
	brokerTopics := b.sub_topics
	b.topic_lock.RUnlock()

	topics := []string{}
	for topic := range brokerTopics {
		topics = append(topics, topic)
	}

	return topics
}
