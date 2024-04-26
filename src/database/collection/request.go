package collection

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"metabox-school-bac-giang-request-service/config"
	"metabox-school-bac-giang-request-service/src/utilities"

	"metabox-school-bac-giang-request-service/mongo"
)

const RequestCollection = "request"

const (
	RequestIndexAccountID = "account_id"
)

var (
	_requestCollection        *RequestMongoCollection
	loadRequestRepositoryOnce sync.Once
)

type RequestMongoCollection struct {
	client         *mongo.MongoDB
	collectionName string
	databaseName   string
	indexed        map[string]bool
}

func LoadRequestCollectionMongo(mongoClient *mongo.MongoDB) (err error) {
	loadRequestRepositoryOnce.Do(func() {
		_requestCollection, err = NewRequestMongoCollection(mongoClient, config.Get().DatabaseName)
	})
	return
}

func Request() *RequestMongoCollection {
	if _requestCollection == nil {
		panic("database: like request collection is not initiated")
	}
	return _requestCollection
}

func NewRequestMongoCollection(client *mongo.MongoDB, databaseName string) (*RequestMongoCollection, error) {
	if client == nil {
		return nil, fmt.Errorf("[NewRequestMongoCollection] client nil pointer")
	}
	repo := &RequestMongoCollection{
		client:         client,
		collectionName: RequestCollection,
		databaseName:   databaseName,
		indexed:        make(map[string]bool),
	}
	repo.SetIndex()
	return repo, nil
}

func (repo *RequestMongoCollection) SetIndex() {
	col := repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)

	indexes := []mongoDriver.IndexModel{
		{
			Keys: bson.M{
				RequestIndexAccountID: 1,
			},
			Options: &options.IndexOptions{
				Name:   utilities.SetString(RequestIndexAccountID),
				Unique: utilities.SetBool(true),
			},
		},
	}

	if !repo.needIndex(col) {
		return
	}

	col.Indexes().CreateMany(context.Background(), indexes)
}

func (repo *RequestMongoCollection) needIndex(col *mongoDriver.Collection) bool {
	keyIndexes := []string{
		RequestIndexAccountID,
	}

	listIndexes, err := col.Indexes().ListSpecifications(context.Background())
	if err != nil {
		return true
	}
	indexed := make([]string, 0)
	for i := 0; i < len(listIndexes); i++ {
		indexed = append(indexed, listIndexes[i].Name)
	}

	for i := 0; i < len(keyIndexes); i++ {
		if !utilities.StringIntArray(keyIndexes[i], indexed) {
			return true
		}
	}

	return false
}

func (repo *RequestMongoCollection) Collection() *mongoDriver.Collection {
	return repo.client.Client().Database(repo.databaseName).Collection(repo.collectionName)
}
