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

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Ack() error {
	p.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(p.queueUrl),
		ReceiptHandle: aws.String(p.receiptHandle),
	})
	return nil
}
