package database

import (
	"context"
	"time"

	"metabox-school-bac-giang-request-service/mongo"

	"metabox-school-bac-giang-request-service/config"
	"metabox-school-bac-giang-request-service/src/database/collection"

	src_const "metabox-school-bac-giang-request-service/src/const"
)

func ConnectDatabse(ctx context.Context) error {
	var mongoClient *mongo.MongoDB
	var err error
	numberRetry := config.Get().NumberRetry
	if numberRetry == 0 {
		numberRetry = src_const.DEFAULTNUMBERRETRY
	}

	for i := 1; i <= config.Get().NumberRetry; i++ {
		mongoClient, err = mongo.NewMongoDBFromUrl(ctx, config.Get().MongoURL, time.Second*10)
		if err != nil {
			if i == config.Get().NumberRetry {
				return err
			}
			time.Sleep(10 * time.Second)
		}

		if mongoClient != nil {
			break
		}
	}

	if err := collection.LoadRequestCollectionMongo(mongoClient); err != nil {
		return err
	}
	return nil
}
