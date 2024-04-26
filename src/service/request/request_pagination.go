package service_request

import (
	"context"
	"fmt"
	"log"
	"metabox-school-bac-giang-request-service/src/database/collection"
	"metabox-school-bac-giang-request-service/src/service"
	"strings"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"

	src_const "metabox-school-bac-giang-request-service/src/const"
	model_request "metabox-school-bac-giang-request-service/src/database/model/request"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type RequestPaginationCommand struct {
	SchoolID string

	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
	OrderBy string                 `json:"order_by"`
	Search  map[string]interface{} `json:"search"`
}

func (c *RequestPaginationCommand) Valid() error {
	if c.Page < 1 {
		c.Page = 1
	}

	if c.Limit < 1 {
		c.Limit = 10
	}

	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}
	return nil
}

func RequestPagination(ctx context.Context, c *RequestPaginationCommand) (total int, results []model_request.Request, err error) {
	log.Println("[service_request.RequestPagination] start")
	defer func() {
		log.Println("[service_request.RequestPagination] end", "data", map[string]interface{}{"command: ": *c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return 0, nil, fmt.Errorf(codeErr)
	}

	objOrderBy := bson.M{}
	if c.OrderBy != "" {
		value := src_const.ASC
		if strings.HasPrefix(c.OrderBy, "-") {
			value = src_const.DESC
			c.OrderBy = strings.TrimPrefix(c.OrderBy, "-")
		}

		objOrderBy = bson.M{c.OrderBy: value}
	}

	//Default order by updated_at | new -> old
	if c.OrderBy == "" {
		objOrderBy = bson.M{"updated_at": src_const.DESC}
	}

	// matchStage := bson.D{{Key: "$match", Value: condition}}

	facectStage := bson.D{{
		Key: "$facet",

		Value: bson.M{
			"rows": bson.A{
				bson.M{"$skip": (c.Page - 1) * c.Limit},
				bson.M{"$limit": c.Limit},
			},
			"total": bson.A{
				bson.M{"$count": "count"},
			},
		},
	}}

	sortStage := bson.D{{Key: "$sort", Value: objOrderBy}}

	pipeline := mongoDriver.Pipeline{
		// matchStage,
		sortStage,
		facectStage,
	}

	cur, err := collection.Request().Collection().Aggregate(ctx, pipeline)
	if err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
		return 0, nil, fmt.Errorf(codeErr)
	}

	var listRequest bson.M
	for cur.Next(ctx) {
		err := cur.Decode(&listRequest)
		if err != nil {
			codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
			service.AddError(ctx, "", "", codeErr)
			return 0, nil, fmt.Errorf(codeErr)
		}
	}

	// Extract the total count and rows from the result
	requests := make([]model_request.Request, 0)

	if len(listRequest["total"].(bson.A)) > 0 {
		total = int(listRequest["total"].(bson.A)[0].(bson.M)["count"].(int32))
		rows := listRequest["rows"].(bson.A)

		for _, rawRequest := range rows {
			requestBSON, err := bson.Marshal(rawRequest)
			if err != nil {
				codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
				service.AddError(ctx, "", "", codeErr)
				return 0, nil, fmt.Errorf(codeErr)
			}

			var request model_request.Request
			err = bson.Unmarshal(requestBSON, &request)
			if err != nil {
				codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
				service.AddError(ctx, "", "", codeErr)
				return 0, nil, fmt.Errorf(codeErr)
			}

			requests = append(requests, request)
		}
	}

	return int(total), requests, nil
}
