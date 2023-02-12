package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	config "event-driven/config"
	database "event-driven/database"
	event "event-driven/event"
	models "event-driven/models"
	preprocessing "event-driven/services/preprocessing"
	utils "event-driven/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rabbitmq/amqp091-go"
)

var (
	broker  event.EventBroker
	db      database.DatabaseObject
	options config.Config
)

func startWithRabbit() {
	options = config.AppConfig
	db = database.NewJsonDatabase("data/documents.json")

	broker = event.NewRabbitEvent(options.RABBIT_URI)

	messages, _ := broker.ConsumeMessage(options.TOPIC_CONSUME_SSML_ID)
	ssmlMessages := messages.(<-chan amqp091.Delivery)

	messages, _ = broker.ConsumeMessage(options.TOPIC_CONSUME_NONSSML_ID)
	nonSsmlMessages := messages.(<-chan amqp091.Delivery)

	forever := make(chan bool)

	go func() {
		for d := range ssmlMessages {
			var qMessage models.QueueMessages
			err := json.Unmarshal(d.Body, &qMessage)
			utils.FailsOnError(err, "Cannot decode queue message")

			document, err := db.GetDocument(qMessage.Key)
			utils.FailsOnError(err, "Cannot get document from database")

			cleanedText := preprocessing.PreprocessSsml(document.Text)
			log.Println(cleanedText)
			document.CleanText = cleanedText

			err = db.UpdateDocument(qMessage.Key, document)
			utils.FailsOnError(err, "Cannot update document to database")

			err = broker.PublishMessage(d.Body, options.TOPIC_PUBLISH_SSML_ID)
			utils.FailsOnError(err, "Cannot publish message to queue")
		}
	}()

	go func() {
		for d := range nonSsmlMessages {
			var qMessage models.QueueMessages
			err := json.Unmarshal(d.Body, &qMessage)
			utils.FailsOnError(err, "Cannot decode queue message")

			document, err := db.GetDocument(qMessage.Key)
			utils.FailsOnError(err, "Cannot get document from database")

			cleanedText := preprocessing.PreprocessNonSsml(document.Text)
			log.Println(cleanedText)
			document.CleanText = cleanedText

			err = db.UpdateDocument(qMessage.Key, document)
			utils.FailsOnError(err, "Cannot update document to database")

			err = broker.PublishMessage(d.Body, options.TOPIC_PUBLISH_NONSSML_ID)
			utils.FailsOnError(err, "Cannot publish message to queue")
		}
	}()

	log.Println("[*] Waiting for messages. Press CTRL+C to cancel")

	<-forever
}

func startWithLambdaSqs(ctx context.Context, events events.SQSEvent) (string, error) {

	broker = event.NewSqsEvent(config.AppConfig.SQS_URL)

	for _, message := range events.Records {

		var qMessage models.QueueMessages
		err := json.Unmarshal([]byte(message.Body), &qMessage)
		utils.FailsOnError(err, "Cannot decode queue message")

		document, err := db.GetDocument(qMessage.Key)
		utils.FailsOnError(err, "Cannot get document from database")

		switch document.SynthesizeType {

		case "ssml":
			cleanedText := preprocessing.PreprocessSsml(document.Text)
			log.Println(cleanedText)
			document.CleanText = cleanedText

			err = db.UpdateDocument(qMessage.Key, document)
			utils.FailsOnError(err, "Cannot update document to database")

			err = broker.PublishMessage([]byte(message.Body), "")
			utils.FailsOnError(err, "Cannot publish message to queue")

			return "Lambda ssml function have been executed", err

		case "non_ssml":
			cleanedText := preprocessing.PreprocessNonSsml(document.Text)
			log.Println(cleanedText)
			document.CleanText = cleanedText

			err = db.UpdateDocument(qMessage.Key, document)
			utils.FailsOnError(err, "Cannot update document to database")

			err = broker.PublishMessage([]byte(message.Body), "")
			utils.FailsOnError(err, "Cannot publish message to queue")

			return "Lambda non_ssml function have been executed", err
		}

	}

	return "Failed to process data", errors.New("synthesize type haven't been found")
}

func main() {
	startWithRabbit()

	if false {
		lambda.Start(startWithLambdaSqs)
	}
}
