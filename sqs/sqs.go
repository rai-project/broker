package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func test() {

	var name string
	var timeout int64
	// Initialize a session that the SDK will use to load configuration,
	// credentials, and region from the shared config file. (~/.aws/config).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a SQS service client.
	svc := sqs.New(sess)

	// Need to convert the queue name into a URL. Make the GetQueueUrl
	// API call to retrieve the URL. This is needed for receiving messages
	// from the queue.
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "AWS.SimpleQueueService.NonExistentQueue" {
			exitErrorf("Unable to find queue %q.", name)
		}
		exitErrorf("Unable to queue %q, %v.", name, err)
	}

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
		exitErrorf("Unable to receive message from queue %q, %v.", name, err)
	}

	fmt.Printf("Received %d messages.\n", len(result.Messages))
	if len(result.Messages) > 0 {
		fmt.Println(result.Messages)
	}
}
