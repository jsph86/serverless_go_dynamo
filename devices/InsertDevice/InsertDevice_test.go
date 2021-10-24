package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedPutItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.PutItemOutput
}

func (d mockedPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &d.Response, nil
}

func TestHandler(t *testing.T) {

	t.Run("Insert New Item To DB", func(t *testing.T) {

		m := mockedPutItem{
			Response: dynamodb.PutItemOutput{},
		}

		d := deps{
			ddb:   m,
			table: "devices",
		}

		e := events.APIGatewayProxyRequest{
			HTTPMethod: "POST",
			Path:       "/insertdevice",
			Body:       `{ "id": "Tobi" }`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		_, err := d.Handler(e)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Insert New Item Without ID Attribute", func(t *testing.T) {
		m := mockedPutItem{
			Response: dynamodb.PutItemOutput{},
		}

		d := deps{
			ddb:   m,
			table: "devices",
		}

		e := events.APIGatewayProxyRequest{
			HTTPMethod: "POST",
			Path:       "/insertdevice",
			Body:       `{ "idd": "Tobi" }`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		res, _ := d.Handler(e)
		if res.StatusCode != 201 {
			t.Fatal("Test Failed : " + res.Body)
		}
	})

}
