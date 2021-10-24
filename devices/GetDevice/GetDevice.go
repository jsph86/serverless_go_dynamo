package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Response events.APIGatewayProxyResponse

type Device struct {
	Id          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

type ResponseObject struct {
	ResponseCode    int
	ResponseContent string
}

func handleGetItem(req events.APIGatewayProxyRequest) (Response, error) {

	responseObject := ResponseObject{0, ""}

	idForSearch := req.QueryStringParameters["id"]

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := "devices"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(idForSearch),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		// Return 500 Internal Server Error
		responseObject = ResponseObject{500, "Internal Server Error"}
	}

	item := Device{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		log.Fatalf("Failed to unmarshal Record, %v", err)
		// Return 500 Internal Server Error
		responseObject = ResponseObject{500, "Internal Server Error"}
	}

	if item.Id == "" {
		// Return 404 Not Found
		responseObject = ResponseObject{404, "Not Found"}
	}

	fmt.Println("item", item)
	out, _ := json.Marshal(item)

	// Return 404 Not Found
	responseObject = ResponseObject{201, string(out)}

	// Return Final Response To ApiGateWay
	return Response{
		StatusCode:      responseObject.ResponseCode,
		IsBase64Encoded: false,
		Body:            responseObject.ResponseContent,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(handleGetItem)
}
