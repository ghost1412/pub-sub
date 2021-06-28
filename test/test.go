package main

import (
	"fmt"
	"log"
	"math/rand"
	pubsub "pubsub"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RandPublisher struct {
}

func (p RandPublisher) publish(topic *pubsub.Topics, topic_name string, message string) {
	for {
		if "" == message {
			message = fmt.Sprint(rand.Float64() * 1000)
		}
		fmt.Println("Sending: ", message)
		topic.Broadcast(message, topic_name)
		time.Sleep(time.Second)
	}
}

// func testTopics() {
// 	topic := pubsub.NewTopic("topic1")
// 	topic.CreateTopic("topic2")
// 	reg_topics := topic.GetTopics()
// 	fmt.Println(reflect.TypeOf(reg_topics))
// 	topics := [3]string{"topic1", "topic2"}

// 	fmt.Println("passed")
// }

func main() {
	topic := pubsub.NewTopic("topic1")
	topic.CreateTopic("topic2")

	p := RandPublisher{}
	subscriber1, err := topic.Register_Sub()
	if err != nil {
		log.Println("Error", err.Error())
	}

	subscriber2, err := topic.Register_Sub()
	if err != nil {
		log.Println("Error", err.Error())
	}
	topic.Subscribe(subscriber1, "topic1")
	topic.Subscribe(subscriber2, "topic2")
	fmt.Println("Subscribers: ", topic.Subscribers("topic1"))
	fmt.Println("Subscribers: ", topic.Subscribers("topic2"))

	topic.CreateTopic("topic3")

	fmt.Println(topic.GetTopics())
	topic.DeleteTopic("topic3")
	fmt.Println(topic.GetTopics())

	ch1 := subscriber1.GetMessages()
	ch2 := subscriber2.GetMessages()

	go p.publish(topic, "topic1", "")
	go p.publish(topic, "topic2", "")

	go subscriber1.Ack(subscriber1.GetID(), ch1)
	go subscriber2.Ack(subscriber2.GetID(), ch2)

	fmt.Scanln()
}
