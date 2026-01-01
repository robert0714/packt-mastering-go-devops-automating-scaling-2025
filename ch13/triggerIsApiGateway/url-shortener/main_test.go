package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name               string
		request            events.APIGatewayProxyRequest
		expectedBody       string
		expectedStatusCode int
		expectError        bool
	}{
		{
			// Test successful URL shortening
			name: "valid URL request",
			request: events.APIGatewayProxyRequest{
				Body: `{"url": "https://example.com/very/long/url"}`,
			},
			expectedBody:       `{"short_url":"https://short.ly/abc123"}`,
			expectedStatusCode: 200,
			expectError:        false,
		},
		{
			// Test with empty body
			name: "empty body",
			request: events.APIGatewayProxyRequest{
				Body: "",
			},
			expectedStatusCode: 400,
			expectError:        true,
		},
		{
			// Test with invalid JSON
			name: "invalid JSON",
			request: events.APIGatewayProxyRequest{
				Body: "invalid json",
			},
			expectedStatusCode: 400,
			expectError:        true,
		},
		{
			// Test with empty URL in request
			name: "empty URL in request",
			request: events.APIGatewayProxyRequest{
				Body: `{"url": ""}`,
			},
			expectedBody:       `{"short_url":"https://short.ly/abc123"}`,
			expectedStatusCode: 200,
			expectError:        false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := handler(context.Background(), testCase.request)

			if testCase.expectError && err == nil {
				t.Errorf("Expected an error, but got nil")
			}

			if !testCase.expectError && err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}

			if response.StatusCode != testCase.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", testCase.expectedStatusCode, response.StatusCode)
			}

			if !testCase.expectError && response.Body != testCase.expectedBody {
				t.Errorf("Expected response body %v, but got %v", testCase.expectedBody, response.Body)
			}
		})
	}
}
