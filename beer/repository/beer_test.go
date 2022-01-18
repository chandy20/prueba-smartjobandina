package repository

import (
	"errors"
	"github.com/chandy20/prueba-smartjobandina/beer/model"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/ory/dockertest"
)

//createBeersTable function to create table beers for testt
func createBeersTable(client *dynamodb.DynamoDB, table string, t *testing.T) {
	_, err := client.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(table),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("active"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("by_active"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("active"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Could not create table: %s", err)
	}
}

func TestBeerRepository_SaveAndFind(t *testing.T) {
	tableBeers := "table_warehouses" + postfix()
	closer, client := dynamodbServerStart(t)
	defer closer()
	createBeersTable(client, tableBeers, t)

	beerToSave := model.Beer{
		ID:       1,
		Name:     "Pilsen",
		Brewery:  "Bavaria",
		Country:  "Colombia",
		Price:    2400,
		Currency: "COP",
	}

	beerRepository := NewBeerRepository(client, tableBeers, logrus.New())

	err := beerRepository.Save(beerToSave)
	if err != nil {
		t.Errorf("error saving berr %v", err)
	}

	beerGot, err := beerRepository.Find(beerToSave.ID)
	if err != nil {
		t.Errorf("error finding beer %v", err)
	}

	if diff := cmp.Diff(beerToSave, beerGot); diff != "" {
		t.Errorf("Error, saved beer is different than expected, (-want,+got)\n%s", diff)
	}

	beerGot, err = beerRepository.Find(2)
	if err != nil {
		t.Errorf("error finding beer %v", err)
	}

	if beerGot.ID > 0 {
		t.Errorf("the test musn't return any beer but return %v", beerGot)
	}

}

func TestBeerRepository_SaveAndList(t *testing.T) {
	tableBeers := "table_warehouses" + postfix()
	closer, client := dynamodbServerStart(t)
	defer closer()
	createBeersTable(client, tableBeers, t)

	beersToSave := []model.Beer{
		{
			ID:       1,
			Name:     "Pilsen",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2400,
			Currency: "COP",
		},
		{
			ID:       2,
			Name:     "Brava",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2000,
			Currency: "COP",
		},
		{
			ID:       3,
			Name:     "Corona",
			Brewery:  "Bavaria",
			Country:  "Mexico",
			Price:    200,
			Currency: "MXN",
		},
		{
			ID:       4,
			Name:     "Budweiser",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2.5,
			Currency: "USD",
		},
		{
			ID:       5,
			Name:     "Leona",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    2000,
			Currency: "COP",
		},
		{
			ID:       6,
			Name:     "Reds",
			Brewery:  "Bavaria",
			Country:  "Colombia",
			Price:    3,
			Currency: "EUR",
		},
	}

	beerRepository := NewBeerRepository(client, tableBeers, logrus.New())

	for _, beer := range beersToSave {
		err := beerRepository.Save(beer)
		if err != nil {
			t.Errorf("error saving berr %v", err)
		}
	}

	beersGot, err := beerRepository.List()
	if err != nil {
		t.Errorf("error listing beer %v", err)
	}

	if len(beersGot) != len(beersToSave) {
		t.Errorf("test must return %v elements but return %v elemens", len(beersGot), len(beersToSave))
	}

}

//postfix function to build a postfix for test tables
func postfix() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// portActive pausa la ejecucion hasta que el socket esta activo o los intentos se agotan
func portActive(network, address string, max int) error {
	for i := 0; i < max; i++ {
		s, err := net.Dial(network, address)
		if err == nil {
			s.Close()
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return errors.New("port is not open")
}

// dynamodbServerStart lanza un servidor dynamodb local para pruebas
func dynamodbServerStart(t *testing.T) (func(), *dynamodb.DynamoDB) {

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "x")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "x")

	dynamodbURL := os.Getenv("DYNAMODB_URL")

	closer := func() {}

	if dynamodbURL == "" {
		pool, err := dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
		// pulls an image, creates a container based on it and runs it
		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository:   "amazon/dynamodb-local",
			Tag:          "latest",
			ExposedPorts: []string{"8000"},
		})
		if err != nil {
			t.Fatalf("Could not start resource: %s", err)
		}
		err = portActive("tcp", resource.GetHostPort("8000/tcp"), 1000)
		if err != nil {
			t.Fatalf("Could not connect resource: %s", resource.GetHostPort("8000/tcp"))
		}
		dynamodbURL = "http://" + resource.GetHostPort("8000/tcp")
		closer = func() {
			if err := pool.Purge(resource); err != nil {
				t.Fatal(err)
			}
		}
	}
	session, err := session.NewSession()
	if err != nil {
		t.Errorf("Error while creating dynamodb local server: %v\n", err)
	}
	client := dynamodb.New(session, &aws.Config{Endpoint: aws.String(dynamodbURL)})
	return closer, client
}
