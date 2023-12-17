package battle

import (
	"context"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	domain "github.com/toledoom/gork/internal/domain/battle"
	"github.com/toledoom/gork/pkg/gork"
)

const tableName = "Battles"

type DynamoStorage struct {
	d *dynamodb.Client
}

func NewDynamoStorage(d *dynamodb.Client) *DynamoStorage {
	return &DynamoStorage{
		d: d,
	}
}

func (br *DynamoStorage) GetByID(id string) (gork.Entity, error) {
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(tableName),
	}

	battleDynamodb, err := br.d.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return nil, err
	}

	b := &domain.Battle{}
	err = attributevalue.UnmarshalMap(battleDynamodb.Item, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (br *DynamoStorage) Add(b gork.Entity) error {
	marshaledItem, err := attributevalue.MarshalMap(b)
	if err != nil {
		return err
	}

	addItemInput := &dynamodb.PutItemInput{
		Item:      marshaledItem,
		TableName: aws.String(tableName),
	}

	_, err = br.d.PutItem(context.TODO(), addItemInput)
	return err
}

func (br *DynamoStorage) Update(e gork.Entity) error {

	b := e.(*domain.Battle)

	updateItemInput := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: b.ID},
		},
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String("SET Player1Score = :Player1Score, Player2Score = :Player2Score"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Player1Score": &types.AttributeValueMemberN{Value: strconv.FormatInt(b.OriginalPlayer1Score, 10)},
			":Player2Score": &types.AttributeValueMemberN{Value: strconv.FormatInt(b.OriginalPlayer2Score, 10)},
		},
	}

	_, err := br.d.UpdateItem(context.TODO(), updateItemInput)

	return err
}

type UowRepository struct {
	uow gork.Worker
}

func NewUowRepository(uow gork.Worker) *UowRepository {
	return &UowRepository{
		uow: uow,
	}
}

func (ur *UowRepository) Add(b *domain.Battle) error {
	return ur.uow.RegisterNew(b)
}

func (ur *UowRepository) Update(b *domain.Battle) error {
	return ur.uow.RegisterDirty(b)
}

func (ur *UowRepository) GetByID(id string) (*domain.Battle, error) {
	entity, err := ur.uow.FetchOne(reflect.TypeOf(&domain.Battle{}), id)
	if err != nil {
		return nil, err
	}
	b := entity.(*domain.Battle)

	return b, nil
}
