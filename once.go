package once

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type manager struct {
	dynamoTableName string

	// TODO auth
}

func New(dynamoTableName string, ops ...option) *manager {
	result := &manager{
		dynamoTableName: dynamoTableName,
	}

	for _, o := range ops {
		o(result)
	}

	return result
}

type option func(m *manager)

func (m *manager) Ensure(id string) {
	svc := dynamodb.New(session.Must(session.NewSession()))

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:              aws.String(m.dynamoTableName),
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
			default:
				fmt.Printf("once发生了没有处理的aws错误: %s", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Printf("once发生了不知道啥错误: %s", err)
		}

		panic(err)
	}

	return
}
