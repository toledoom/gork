package player

import (
	"context"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	domain "github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/gork"
)

const tableName = "Players"

type DynamoStorage struct {
	d *dynamodb.Client
}

func NewDynamoStorage(d *dynamodb.Client) *DynamoStorage {
	return &DynamoStorage{
		d: d,
	}
}

func (pr *DynamoStorage) GetByID(id string) (gork.Entity, error) {
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(tableName),
	}

	playerDynamodb, err := pr.d.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return nil, err
	}

	p := &domain.Player{}
	err = attributevalue.UnmarshalMap(playerDynamodb.Item, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (pr *DynamoStorage) Add(p gork.Entity) error {
	marshaledItem, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	addItemInput := &dynamodb.PutItemInput{
		Item:      marshaledItem,
		TableName: aws.String(tableName),
	}

	_, err = pr.d.PutItem(context.TODO(), addItemInput)
	return err
}

func (br *DynamoStorage) Update(e gork.Entity) error {
	p := e.(*domain.Player)
	updateItemInput := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: p.ID},
		},
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String("SET Score = :Score"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Score": &types.AttributeValueMemberN{Value: strconv.FormatInt(p.Score, 10)},
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

func (ur *UowRepository) Add(p *domain.Player) error {
	return ur.uow.RegisterNew(p)
}

func (ur *UowRepository) Update(p *domain.Player) error {
	return ur.uow.RegisterDirty(p)
}

func (ur *UowRepository) Delete(p *domain.Player) error {
	return ur.uow.RegisterDeleted(p)
}

func (ur *UowRepository) GetByID(id string) (*domain.Player, error) {
	entity, err := ur.uow.FetchOne(reflect.TypeOf(&domain.Player{}), id)
	if err != nil {
		return nil, err
	}
	b := entity.(*domain.Player)

	return b, nil
}
