package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Response events.APIGatewayProxyResponse
type deps struct {
	ddb   dynamodbiface.DynamoDBAPI
	table string
}
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

func ResponseToGateway(response ResponseObject, errorToThrow error) (Response, error) {
	return Response{
		StatusCode:      response.ResponseCode,
		IsBase64Encoded: false,
		Body:            response.ResponseContent,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, errorToThrow
}

func (d *deps) Handler(req events.APIGatewayProxyRequest) (Response, error) {

	responseObject := ResponseObject{0, ""}

	var errorToThrow error

	if req.Body != "" {

		var newDevice Device

		//Convert Delivered Object Via ApiGateway to Device Object
		errorWhileCOnverting := json.Unmarshal([]byte(req.Body), &newDevice)

		if errorWhileCOnverting != nil || newDevice.Id == "" {
			responseObject = ResponseObject{500, "required Fields is not available"}
			return ResponseToGateway(responseObject, errorToThrow)
		}

		// Convert Device Object to Valid Dynamodb Object
		av, err := dynamodbattribute.MarshalMap(newDevice)

		if err != nil {
			log.Fatalf("Got Error Creating Object For Saving : %s", err)
			responseObject = ResponseObject{400, "Bad Request"}
			return ResponseToGateway(responseObject, errorToThrow)
		}

		tableName := d.table

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		// Send Data to DynamoDB
		_, err = d.ddb.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %s", err)

			responseObject = ResponseObject{400, "Bad Request"}
			return ResponseToGateway(responseObject, errorToThrow)

		} else {
			responseObject = ResponseObject{201, "Created."}
			return ResponseToGateway(responseObject, errorToThrow)
		}

	} else {
		responseObject = ResponseObject{400, "Bad Request"}
		return ResponseToGateway(responseObject, errorToThrow)
	}

}

func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	d := deps{
		svc,
		"devices",
	}

	lambda.Start(d.Handler)
}
