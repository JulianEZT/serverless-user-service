package users

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	pkPrefix = "USER#"
	skValue  = "PROFILE"
)

// DynamoRepo implements UserRepository with DynamoDB.
type DynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoRepo returns a DynamoRepo.
func NewDynamoRepo(client *dynamodb.Client, tableName string) *DynamoRepo {
	return &DynamoRepo{client: client, tableName: tableName}
}

// dynamoUser is the stored item shape (pk, sk, and attributes).
type dynamoUser struct {
	PK        string `dynamodbav:"pk"`
	SK        string `dynamodbav:"sk"`
	ID        string `dynamodbav:"id"`
	Email     string `dynamodbav:"email"`
	Name      string `dynamodbav:"name"`
	CreatedAt string `dynamodbav:"createdAt"`
	CreatedBy string `dynamodbav:"createdBy"`
}

func toDynamo(u *User) dynamoUser {
	return dynamoUser{
		PK:        pkPrefix + u.ID,
		SK:        skValue,
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
	}
}

func (d *DynamoRepo) Put(ctx context.Context, u *User) error {
	item, err := attributevalue.MarshalMap(toDynamo(u))
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
		ConditionExpression: ptr("attribute_not_exists(pk)"),
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *DynamoRepo) GetByID(ctx context.Context, id string) (*User, error) {
	pk := pkPrefix + id
	out, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: skValue},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	var du dynamoUser
	if err := attributevalue.UnmarshalMap(out.Item, &du); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &User{
		ID:        du.ID,
		Email:     du.Email,
		Name:      du.Name,
		CreatedAt: du.CreatedAt,
		CreatedBy: du.CreatedBy,
	}, nil
}

func ptr(s string) *string { return &s }
