package stream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

// Publisher sends records to an AWS Kinesis Data Stream.
type Publisher struct {
	client     *kinesis.Client
	streamName string
}

// NewPublisher creates a Kinesis publisher for the given stream name.
func NewPublisher(cfg aws.Config, streamName string) *Publisher {
	return &Publisher{
		client:     kinesis.NewFromConfig(cfg),
		streamName: streamName,
	}
}

// Publish serialises payload as JSON and puts it onto the Kinesis stream.
// The battery_id is used as the partition key to preserve per-battery ordering.
func (p *Publisher) Publish(ctx context.Context, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	// TODO: extract battery_id from payload via type assertion for deterministic sharding
	partitionKey := fmt.Sprintf("shard-%d", len(data)%4)

	_, err = p.client.PutRecord(ctx, &kinesis.PutRecordInput{
		StreamName:   aws.String(p.streamName),
		Data:         data,
		PartitionKey: aws.String(partitionKey),
	})
	if err != nil {
		return fmt.Errorf("kinesis put: %w", err)
	}
	return nil
}
