package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go-lambda-ex/pkg/handlers"
	"log"
	"os"
)

const dynamoDbTableName = "Users"

var (
	dynamoClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Printf("Error in initializing AWS session %s \n", err.Error())
		return
	}
	dynamoClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, dynamoDbTableName, dynamoClient)
	case "POST":
		return handlers.CreateUser(req, dynamoDbTableName, dynamoClient)
	case "PUT":
		return handlers.UpdateUser(req, dynamoDbTableName, dynamoClient)
	case "DELETE":
		return handlers.DeleteUser(req, dynamoDbTableName, dynamoClient)
	default:
		return handlers.UnhandledMethod()
	}
}
