package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rai-project/broker"
)

type publication struct {
	svc           *sqs.SQS
	topic         string
	m             *broker.Message
	nm            *sqs.Message
	queueUrl      string
	receiptHandle string
	opts          broker.PublishOptions
}

// Topic ...
func (p *publication) Topic() string {
	return p.topic
}

// Message ...
func (p *publication) Message() *broker.Message {
	return p.m
}

// Ack ...
func (p *publication) Ack() error {
	_, err := p.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(p.queueUrl),
		ReceiptHandle: aws.String(p.receiptHandle),
	})
	return err
}
