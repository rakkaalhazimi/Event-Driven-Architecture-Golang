package test

import (
	"testing"

	config "event-driven/config"
	event "event-driven/event"
)

var options config.Config

func init() {
	options = config.AppConfig
}

func TestPublishMessage(t *testing.T) {
	broker := event.NewRabbitEvent(options.RABBIT_URI)
	sampleMessage := []byte(`{"key":"0"}`)
	t.Log(options.TOPIC_CONSUME_NONSSML_ID)
	err := broker.PublishMessage(sampleMessage, options.TOPIC_CONSUME_NONSSML_ID)
	if err != nil {
		t.Fail()
	}
}
