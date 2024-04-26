package service_request

import (
	"context"
	"fmt"
	"log"
	src_const "metabox-school-bac-giang-request-service/src/const"
	"metabox-school-bac-giang-request-service/src/database/collection"
	model_request "metabox-school-bac-giang-request-service/src/database/model/request"
	"metabox-school-bac-giang-request-service/src/service"

	"go.mongodb.org/mongo-driver/bson"
)

type RequestDeleteCommand struct {
	RequestID string
}

func (c *RequestDeleteCommand) Valid() error {
	if c.RequestID == "" {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func RequestDelete(ctx context.Context, c *RequestDeleteCommand) (result *model_request.Request, err error) {
	log.Println("[service_request.RequestDelete] start")
	defer func() {
		log.Println("[service_request.RequestDelete] end", "data", map[string]interface{}{"command: ": *c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	question := &model_request.Request{}
	err = collection.Request().Collection().FindOne(ctx, bson.M{"_id": c.RequestID}).Decode(question)
	if err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	_, err = collection.Request().Collection().DeleteOne(ctx, bson.M{"_id": c.RequestID})
	if err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, err
	}

	return question, nil
}
