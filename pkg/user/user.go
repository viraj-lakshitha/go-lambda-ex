package user

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go-lambda-ex/pkg/utils"
	"log"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item error"
	ErrorUserAlreadyExists       = "user.User already exists"
	ErrorUserDoesNotExists       = "user.User does not exist"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		log.Printf("Error in fetching user with email: %s \n", email)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		log.Printf("Error in unmarshaling user with email: %s \n", email)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchUsers(tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := dynamoClient.Scan(input)
	if err != nil {
		log.Println("Error in fetching users")
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		log.Println("Error in unmarshalling users")
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return items, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		log.Println("Error in unable ot unmarshal request body")
		return nil, errors.New(ErrorInvalidUserData)
	}

	// Validate email address
	if !utils.IsEmailValid(u.Email) {
		log.Println("Invalid email address in request body")
		return nil, errors.New(ErrorInvalidEmail)
	}

	// Check for existing records
	currentUser, _ := FetchUser(u.Email, tableName, dynamoClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		log.Printf("User already exist with email %s /n", u.Email)
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	// Save record
	userStr, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		log.Println("Error in marshaling request")
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      userStr,
		TableName: aws.String(tableName),
	}
	_, err = dynamoClient.PutItem(input)
	if err != nil {
		log.Println("Unable to create new record")
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, err
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		log.Println("Unable to unmarshal request")
		return nil, errors.New(ErrorInvalidEmail)
	}

	// Check if user exists
	currentUser, _ := FetchUser(u.Email, tableName, dynamoClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		log.Println("Unable to fetch existing user")
		return nil, errors.New(ErrorUserDoesNotExists)
	}

	// Save record
	userStr, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		log.Println("Error in marshaling request")
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      userStr,
		TableName: aws.String(tableName),
	}
	_, err = dynamoClient.PutItem(input)
	if err != nil {
		log.Println("Unable to create new record")
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, err
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]

	// Validate email address
	if !utils.IsEmailValid(email) {
		log.Printf("Invalid email address %s \n", email)
		return errors.New(ErrorInvalidEmail)
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		log.Println("Unable to delete record")
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}
