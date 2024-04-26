package service_request

import (
	"context"
	"fmt"
	"log"
	"metabox-school-bac-giang-request-service/src/database/collection"
	"metabox-school-bac-giang-request-service/src/service"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	src_const "metabox-school-bac-giang-request-service/src/const"
	model_request "metabox-school-bac-giang-request-service/src/database/model/request"
)

type RequestAddCommand struct {
	Username string

	Note  string
	Phone string
}

func (c *RequestAddCommand) Valid() error {
	if c.Username == "" {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func RequestAdd(ctx context.Context, c *RequestAddCommand) (result *model_request.Request, err error) {
	log.Println("[service_request.RequestAdd] start")
	defer func() {
		log.Println("[service_request.RequestAdd] end", "data", map[string]interface{}{"command: ": *c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	result = &model_request.Request{
		ID: primitive.NewObjectID().Hex(),

		Username: c.Username,
		Phone:    c.Phone,
		Note:     c.Note,

		Deleted:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = collection.Request().Collection().InsertOne(ctx, result)
	if err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}
	return result, nil
}
