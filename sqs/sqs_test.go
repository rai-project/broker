package sqs

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	raiaws "github.com/rai-project/aws"
	"github.com/rai-project/config"

	"github.com/stretchr/testify/assert"
)

func test() error {

	// Initialize a session that the SDK will use to load configuration,
	// credentials, and region from the shared config file. (~/.aws/config).
	sess := session.Must(raiaws.NewSession())

	// Create a SQS service client.
	svc := sqs.New(sess)

	name := "rai"

	// Need to convert the queue name into a URL. Make the GetQueueUrl
	// API call to retrieve the URL. This is needed for receiving messages
	// from the queue.
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "AWS.SimpleQueueService.NonExistentQueue" {
			return errors.Errorf("Unable to find queue %q.", name)
		}
		return errors.Errorf("Unable to queue %q, %v.", name, err)
	}

	timeout := int64(1) // second

	// Send a message from the SQS queue with long polling enabled.
	go func() {
		result, err := svc.SendMessage(&sqs.SendMessageInput{
			QueueUrl:    resultURL.QueueUrl,
			MessageBody: aws.String("fdsa"),
		})
		if err != nil {
			pp.Println("sending error = ", err)
			return
		}
		pp.Println("Message", result.MessageId, "was sent")
	}()

	// Receive a message from the SQS queue with long polling enabled.
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: resultURL.QueueUrl,
		AttributeNames: aws.StringSlice([]string{
			"SentTimestamp",
		}),
		MaxNumberOfMessages: aws.Int64(1),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
		WaitTimeSeconds: aws.Int64(timeout),
	})
	if err != nil {
		return errors.Errorf("Unable to receive message from queue %q, %v.", name, err)
	}

	fmt.Printf("Received %d messages.\n", len(result.Messages))
	if len(result.Messages) > 0 {
		fmt.Println(result.Messages)
	}
	for _, msg := range result.Messages {
		_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      resultURL.QueueUrl,
			ReceiptHandle: msg.ReceiptHandle,
		})
		pp.Println(err)
	}
	return nil
}

// TestSQS ...
func TestSQS(t *testing.T) {
	config.Init()
	assert.NoError(t, test())
}
