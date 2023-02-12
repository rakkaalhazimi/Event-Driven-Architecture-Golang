package event

// Event Policy
type EventBroker interface {
	ConsumeMessage(topic string) (interface{}, error)
	PublishMessage(message interface{}, topic string) error
}
