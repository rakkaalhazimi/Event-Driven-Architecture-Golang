package event

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	utils "event-driven/utils"
)

type SQSBroker struct {
	url     string
	service *sqs.SQS
}

func (e *SQSBroker) ConsumeMessage(topic string) (interface{}, error) {
	err := errors.New("SQS Consume Message Haven't Been Implemented")
	return nil, err
}

func (e *SQSBroker) PublishMessage(message interface{}, topic string) error {

	_, err := e.service.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message.(string)),
		QueueUrl:     &e.url,
	})
	utils.FailsOnError(err, "Cannot publish message to sqs service")

	return err
}

func NewSqsEvent(url string) *SQSBroker {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,

		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	}))
	svc := sqs.New(sess)

	return &SQSBroker{
		url:     url,
		service: svc,
	}
}
