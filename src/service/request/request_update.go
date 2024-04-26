package service_request

import (
	"context"
	"fmt"
	"log"
	src_const "metabox-school-bac-giang-request-service/src/const"
	"metabox-school-bac-giang-request-service/src/database/collection"
	model_request "metabox-school-bac-giang-request-service/src/database/model/request"
	"metabox-school-bac-giang-request-service/src/service"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type RequestUpdateCommand struct {
	RequestID string

	Note  *string
	Phone *string
}

func (c *RequestUpdateCommand) Valid() error {
	if c.RequestID == "" {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func RequestUpdate(ctx context.Context, c *RequestUpdateCommand) (result *model_request.Request, err error) {
	log.Println("[service_request.RequestUpdate] start")
	defer func() {
		log.Println("[service_request.RequestUpdate] end", "data", map[string]interface{}{"command: ": *c}, "error", err)
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

	updated := make(map[string]interface{})

	if c.Note != nil {
		updated["note"] = *c.Note
	}

	if c.Phone != nil {
		updated["phone"] = *c.Phone
	}

	updated["updated_at"] = time.Now()
	_, err = collection.Request().Collection().UpdateByID(ctx, c.RequestID, bson.M{"$set": updated})
	if err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
		service.AddError(ctx, "", "", codeErr)
		return nil, err
	}

	return question, nil
}
