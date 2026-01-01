package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name          string
		s3Event       events.S3Event
		expectedError error
	}{
		{
			// mock a S3 event with a single object
			name: "single object upload",
			s3Event: events.S3Event{
				Records: []events.S3EventRecord{
					{
						S3: events.S3Entity{
							Bucket: events.S3Bucket{
								Name: "test-bucket",
							},
							Object: events.S3Object{
								Key:  "test-file.json",
								Size: 1024,
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			// mock a S3 event with multiple objects
			name: "multiple objects upload",
			s3Event: events.S3Event{
				Records: []events.S3EventRecord{
					{
						S3: events.S3Entity{
							Bucket: events.S3Bucket{
								Name: "test-bucket",
							},
							Object: events.S3Object{
								Key:  "file1.json",
								Size: 512,
							},
						},
					},
					{
						S3: events.S3Entity{
							Bucket: events.S3Bucket{
								Name: "test-bucket",
							},
							Object: events.S3Object{
								Key:  "file2.json",
								Size: 2048,
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			// mock a S3 event with no records
			name: "empty records",
			s3Event: events.S3Event{
				Records: []events.S3EventRecord{},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			err := handler(ctx, testCase.s3Event)
			if err != testCase.expectedError {
				t.Errorf("Expected error %v, but got %v", testCase.expectedError, err)
			}
		})
	}
}
