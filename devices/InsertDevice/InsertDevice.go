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

// Handler is our lambda handler invoked by the `lambda.Start` function call
func (d *deps) Handler(req events.APIGatewayProxyRequest) (Response, error) {
	//	var buf bytes.Buffer

	responseObject := ResponseObject{0, ""}

	var errorToThrow error

	if req.Body != "" {

		var newDevice Device

		//Convert Delivered Object Via ApiGateway to Device Object
		erro := json.Unmarshal([]byte(req.Body), &newDevice)

		if erro != nil || newDevice.Id == "" {

			// Return 400 Id is not available
			//	errorToThrow = errors.New("required Fields is not available")
			responseObject = ResponseObject{500, "required Fields is not available"}
			return Response{
				StatusCode:      responseObject.ResponseCode,
				IsBase64Encoded: false,
				Body:            responseObject.ResponseContent,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, errorToThrow
		}

		// Convert Device Object to Valid Dynamodb Object
		av, err := dynamodbattribute.MarshalMap(newDevice)

		if err != nil {
			log.Fatalf("Got Error Creating Object For Saving : %s", err)
			// Return 400 Id is not available
			//	errorToThrow = errors.New("Bad Request")
			responseObject = ResponseObject{400, "Bad Request"}
			return Response{
				StatusCode:      responseObject.ResponseCode,
				IsBase64Encoded: false,
				Body:            responseObject.ResponseContent,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, errorToThrow
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

			// Return 500 Internal Server Error
			//	errorToThrow = errors.New("Bad Request")
			responseObject = ResponseObject{400, "Bad Request"}
			return Response{
				StatusCode:      responseObject.ResponseCode,
				IsBase64Encoded: false,
				Body:            responseObject.ResponseContent,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, errorToThrow

		} else {
			// Return 500 Internal Server Error
			responseObject = ResponseObject{201, "Created."}
			return Response{
				StatusCode:      responseObject.ResponseCode,
				IsBase64Encoded: false,
				Body:            responseObject.ResponseContent,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, errorToThrow
		}

	} else {
		// Return 400 Id is not available
		//	errorToThrow = errors.New("bad Request")
		responseObject = ResponseObject{400, "Bad Request"}
		return Response{
			StatusCode:      responseObject.ResponseCode,
			IsBase64Encoded: false,
			Body:            responseObject.ResponseContent,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, errorToThrow
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
