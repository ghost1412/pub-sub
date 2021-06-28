# pub-sub


The PubSub system should support below methods:
```
    CreateTopic(topicID)
    DeleteTopic(TopicID)
    AddSubscription(topicID,SubscriptionID); Creates and adds subscription with id SubscriptionID to topicName.
    DeleteSubscription(SubscriptionID)
    Publish(topicID, message); publishes the message on given topic
    Subscribe(SubscriptionID, SubscriberFunc); SubscriberFunc is the subscriber which is executed for each message of subscription.
    UnSubscribe(SubscriptionID)
    Ack(SubscriptionID, MessageID); Called by Subscriber to intimate the Subscription that the message has been received and processed
```

to check functionality try ```go run test/test.go```

Classes:

- Topics - can add/remove topics and subscribers
- Messages
- Publisher
- Subscriber 
