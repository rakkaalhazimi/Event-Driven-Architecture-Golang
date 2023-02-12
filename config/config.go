package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BROKER_SERVICE           string
	RABBIT_URI               string
	SQS_URL                  string
	TOPIC_PUBLISH_NONSSML_ID string
	TOPIC_PUBLISH_SSML_ID    string
	TOPIC_CONSUME_NONSSML_ID string
	TOPIC_CONSUME_SSML_ID    string
}

var (
	AppConfig Config
)

func init() {
	godotenv.Load()

	AppConfig = Config{
		BROKER_SERVICE:           os.Getenv("BROKER_SERVICE"),
		RABBIT_URI:               os.Getenv("RABBIT_URI"),
		SQS_URL:                  os.Getenv("SQS_URL"),
		TOPIC_PUBLISH_NONSSML_ID: os.Getenv("TOPIC_PUBLISH_NONSSML_ID"),
		TOPIC_PUBLISH_SSML_ID:    os.Getenv("TOPIC_PUBLISH_SSML_ID"),
		TOPIC_CONSUME_NONSSML_ID: os.Getenv("TOPIC_CONSUME_NONSSML_ID"),
		TOPIC_CONSUME_SSML_ID:    os.Getenv("TOPIC_CONSUME_SSML_ID"),
	}

	fmt.Println("Configuration loaded")
}
