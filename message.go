package pubsub

type Topic string

type creationTime int64

type Message struct {
	topic_name    Topic
	data          interface{}
	creation_time creationTime
}

func (m *Message) GetTopicName() Topic {
	return m.topic_name
}

func (m *Message) GetData() interface{} {
	return m.data
}

func (m *Message) GetCreationTime() creationTime {
	return m.creation_time
}
