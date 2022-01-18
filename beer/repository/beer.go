package repository

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"github.com/sirupsen/logrus"
)

//BeerRepository main struct for repository
type BeerRepository struct {
	client     *dynamodb.DynamoDB
	tableBeers string
	logger     *logrus.Logger
}

//Find method to search a beer
func (b *BeerRepository) Find(ID int) (model.Beer, error) {
	out, err := b.client.Query(&dynamodb.QueryInput{
		TableName: aws.String(b.tableBeers),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(strconv.Itoa(ID)),
					},
				},
			},
		},
	})
	if err != nil {
		return model.Beer{}, err
	}

	if len(out.Items) == 0 {
		return model.Beer{}, nil
	}

	beers, err := b.hydrate(out.Items)
	if err != nil {
		return model.Beer{}, err
	}
	return beers[0], nil
}

//Save method to save a beer
func (b *BeerRepository) Save(beer model.Beer) error {
	logger := b.logger.WithField("model", beer)
	logger.Info("saving beer")
	item := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(strconv.Itoa(beer.ID)),
			},
			"name": {
				S: aws.String(beer.Name),
			},
			"brewery": {
				S: aws.String(beer.Brewery),
			},
			"country": {
				S: aws.String(beer.Country),
			},
			"price": {
				N: aws.String(fmt.Sprintf("%f", beer.Price)),
			},
			"currency": {
				S: aws.String(beer.Currency),
			},
			"active": {
				N: aws.String("1"),
			},
		},
		TableName:           aws.String(b.tableBeers),
		ConditionExpression: aws.String("attribute_not_exists(#id)"),
		ExpressionAttributeNames: map[string]*string{
			"#id": aws.String("id"),
		},
	}
	_, err := b.client.PutItem(item)

	return err
}

//hydrate method to populate []model.Beer from dynamodb.AttributeValue
func (b *BeerRepository) hydrate(items []map[string]*dynamodb.AttributeValue) ([]model.Beer, error) {
	var beers = make([]model.Beer, len(items))
	for i, item := range items {
		value := item["id"].S
		ID, err := strconv.Atoi(*value)
		if err != nil {
			return []model.Beer{}, err
		}
		beers[i].ID = ID

		if v, ok := item["name"]; ok {
			beers[i].Name = *v.S
		}

		if v, ok := item["brewery"]; ok {
			beers[i].Brewery = *v.S
		}

		if v, ok := item["country"]; ok {
			beers[i].Country = *v.S
		}

		if v, ok := item["price"]; ok {
			price, err := strconv.ParseFloat(*v.N, 64)
			if err != nil {
				return []model.Beer{}, err
			}
			beers[i].Price = price
		}

		if v, ok := item["currency"]; ok {
			beers[i].Currency = *v.S
		}
	}
	return beers, nil
}

//List method to list all beers in database
func (b *BeerRepository) List() ([]model.Beer, error) {
	logger := b.logger
	logger.Info("beginning of list beers")
	out, err := b.client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(b.tableBeers),
		IndexName:              aws.String("by_active"),
		KeyConditionExpression: aws.String("active = :active"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":active": {
				N: aws.String("1"),
			},
		},
	})
	if err != nil {
		logger.WithError(err).Error("error listing beers")
		return []model.Beer{}, err
	}
	if len(out.Items) == 0 {
		logger.Info("no beers found")
		return []model.Beer{}, nil
	}

	tmp := out.Items
	lastEvaluatedKey := out.LastEvaluatedKey
	for len(lastEvaluatedKey) != 0 {
		out, err = b.client.Query(&dynamodb.QueryInput{
			TableName:              aws.String(b.tableBeers),
			IndexName:              aws.String("by_active"),
			KeyConditionExpression: aws.String("active = :active"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":active": {
					S: aws.String("1"),
				},
			},
			ExclusiveStartKey: lastEvaluatedKey,
		})
		lastEvaluatedKey = out.LastEvaluatedKey
		if err != nil {
			logger.WithError(err).Error("an error occurred reading another page")
			return []model.Beer{}, err
		}
		tmp = append(tmp, out.Items...)
	}

	beers, err := b.hydrate(tmp)
	if err != nil {
		logger.WithError(err).Error("an error occurred reading another page")
		return []model.Beer{}, err
	}

	return beers, nil
}

//NewBeerRepository construct for repository
func NewBeerRepository(
	client *dynamodb.DynamoDB,
	tableBeers string,
	logger *logrus.Logger,
) *BeerRepository {
	return &BeerRepository{
		client:     client,
		tableBeers: tableBeers,
		logger:     logger,
	}
}
