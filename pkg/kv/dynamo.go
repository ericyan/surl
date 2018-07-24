package kv

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// A DynamoStore is a key-value store backed by a Amazon DynamoDB table.
type DynamoStore struct {
	TableName string
	svc       *dynamodb.DynamoDB
}

const (
	keyAttrName   = "DS_KEY"
	valueAttrName = "DS_VALUE"
)

var (
	primaryKeyDefinition = dynamodb.AttributeDefinition{
		AttributeName: aws.String(keyAttrName),
		AttributeType: aws.String("B"), // B means Binary
	}

	tableKeySchema = dynamodb.KeySchemaElement{
		AttributeName: aws.String(keyAttrName),
		KeyType:       aws.String("HASH"), // HASH attribute means partition key
	}
)

// NewDynamoStore returns a new DynamoStore. If the given table does not
// exist, it will try to create it and wait until the status of the new
// table becames ACTIVE.
func NewDynamoStore(endpoint, tableName string) (*DynamoStore, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	// make sure the table exists
	_, err = svc.CreateTable(&dynamodb.CreateTableInput{
		TableName:            aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{&primaryKeyDefinition},
		KeySchema:            []*dynamodb.KeySchemaElement{&tableKeySchema},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(25),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceInUseException {
			// table already exists, check for validity
			t, err := svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
			if err != nil {
				return nil, err
			}

			if attr := t.Table.AttributeDefinitions[0]; *attr.AttributeName != *primaryKeyDefinition.AttributeName ||
				*attr.AttributeType != *primaryKeyDefinition.AttributeType {
				return nil, errors.New("dynamostore: table exists, but with invalid primary key")
			}

			if schema := t.Table.KeySchema[0]; *schema.AttributeName != *tableKeySchema.AttributeName ||
				*schema.KeyType != *tableKeySchema.KeyType {
				return nil, errors.New("dynamostore: table exists, but with incompatible schema")
			}
		} else {
			return nil, err
		}
	}

	// make sure table is active, check one per 300ms for 3 seconds max
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	cancel := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		cancel <- false
	}()
	for {
		var isActive bool
		select {
		case <-cancel:
			return nil, errors.New("dynamodb: timeout when waiting for table become active")
		case <-ticker.C:
			t, err := svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
			if err != nil {
				return nil, err
			}

			isActive = *t.Table.TableStatus == "ACTIVE"
		}

		if isActive {
			break
		}
	}

	return &DynamoStore{tableName, svc}, nil
}

func (s *DynamoStore) Get(key []byte) ([]byte, error) {
	res, err := s.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyAttrName: new(dynamodb.AttributeValue).SetB(key),
		},
		TableName: aws.String(s.TableName),
	})
	if err != nil {
		return nil, err
	}

	attr, ok := res.Item[valueAttrName]
	if !ok {
		return nil, errors.New("dynamostore: value attribute not found")
	}
	return attr.B, nil
}

func (s *DynamoStore) Put(key, value []byte) error {
	_, err := s.svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			keyAttrName:   new(dynamodb.AttributeValue).SetB(key),
			valueAttrName: new(dynamodb.AttributeValue).SetB(value),
		},
		TableName: aws.String(s.TableName),
	})

	return err
}
