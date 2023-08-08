package clients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	common "github.com/jplindgren/rpg-vault/internal"
)

type DynamoDbClientWrapper struct {
	*dynamodb.Client
}

func (c *DynamoDbClientWrapper) PutWrapper(tableName string, item interface{}, conditionExp *string) (*dynamodb.PutItemOutput, error) {
	av, marshalErr := attributevalue.MarshalMap(item)
	if marshalErr != nil {
		return &dynamodb.PutItemOutput{}, marshalErr
	}

	putItemRes, putItemErr := c.PutItem(context.TODO(), &dynamodb.PutItemInput{
		//TableName:           aws.String(config.PrimaryTableName),
		TableName:           aws.String(tableName),
		Item:                av,
		ConditionExpression: conditionExp,
	})
	if putItemErr != nil {
		return &dynamodb.PutItemOutput{}, putItemErr
	}

	return putItemRes, nil
}

func (c *DynamoDbClientWrapper) GetWrapper(tableName string, key interface{}, resultItem interface{}) (*dynamodb.GetItemOutput, error) {
	av, marshalErr := attributevalue.MarshalMap(key)
	if marshalErr != nil {
		return &dynamodb.GetItemOutput{}, marshalErr
	}

	getItemRes, getItemErr := c.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       av,
	})

	if getItemErr != nil {
		return &dynamodb.GetItemOutput{}, getItemErr
	}

	if getItemRes.Item == nil {
		return &dynamodb.GetItemOutput{}, common.ErrorRecordNotFound
	}

	unmarshalErr := attributevalue.UnmarshalMap(getItemRes.Item, resultItem)
	if unmarshalErr != nil {
		return &dynamodb.GetItemOutput{}, unmarshalErr
	}

	return getItemRes, nil
}

func (c *DynamoDbClientWrapper) QueryWrapper(tableName string, keyCondition expression.KeyConditionBuilder, resultArr interface{}) (interface{}, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, err
	}

	return c.QueryWithExpressionWrapper(tableName, expr, resultArr)
}

func (c *DynamoDbClientWrapper) QueryWithExpressionWrapper(tableName string, expr expression.Expression, resultArr interface{}) ([]map[string]types.AttributeValue, error) {
	var response *dynamodb.QueryOutput

	response, err := c.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &resultArr)
	if err != nil {
		return nil, err
	}
	return response.Items, err
}

func (c *DynamoDbClientWrapper) UpdateWrapper(tableName string, key interface{}, update expression.UpdateBuilder) (*dynamodb.UpdateItemOutput, error) {
	av, marshalErr := attributevalue.MarshalMap(key)
	if marshalErr != nil {
		return &dynamodb.UpdateItemOutput{}, marshalErr
	}

	expr, builderErr := expression.NewBuilder().WithUpdate(update).Build()
	if builderErr != nil {
		return &dynamodb.UpdateItemOutput{}, builderErr
	}

	updateItemRes, updateItemErr := c.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		//TableName:           aws.String(config.PrimaryTableName),
		TableName:                 aws.String(tableName),
		Key:                       av,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if updateItemErr != nil {
		return &dynamodb.UpdateItemOutput{}, updateItemErr
	}

	return updateItemRes, nil
}

func (c *DynamoDbClientWrapper) DeleteWrapper(tableName string, key interface{}) (*dynamodb.DeleteItemOutput, error) {
	av, marshalErr := attributevalue.MarshalMap(key)
	if marshalErr != nil {
		return &dynamodb.DeleteItemOutput{}, marshalErr
	}

	deleteItemRes, deleteItemErr := c.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		//TableName:           aws.String(config.PrimaryTableName),
		TableName: aws.String(tableName),
		Key:       av,
	})
	if deleteItemErr != nil {
		return &dynamodb.DeleteItemOutput{}, deleteItemErr
	}

	return deleteItemRes, nil
}

type CharacterKey struct {
	WorldId string
	Id      string
}

func (c *DynamoDbClientWrapper) BatchDeleteWrapper(tableName string, keys []map[string]string) (*dynamodb.BatchWriteItemOutput, error) {
	var wr []types.WriteRequest
	for _, key := range keys {
		parsedKey, err := attributevalue.MarshalMap(key)
		if err != nil {
			return nil, err
		}

		wr = append(
			wr,
			types.WriteRequest{DeleteRequest: &types.DeleteRequest{Key: parsedKey}},
		)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{tableName: wr}}

	deleteItemRes, deleteItemErr := c.BatchWriteItem(context.TODO(), input)
	if deleteItemErr != nil {
		return nil, deleteItemErr
	}

	return deleteItemRes, nil
}
