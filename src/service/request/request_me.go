package service_request

import (
	"context"
	"fmt"
	"log"
	"metabox-school-bac-giang-request-service/src/service"

	src_const "metabox-school-bac-giang-request-service/src/const"
	"metabox-school-bac-giang-request-service/src/database/collection"
	model_request "metabox-school-bac-giang-request-service/src/database/model/request"
)

type RequestMeCommand struct {
	Username string
}

func (c *RequestMeCommand) Valid() error {
	if c.Username == "" {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		return fmt.Errorf(codeErr)
	}

	return nil
}

func RequestMe(ctx context.Context, c *RequestMeCommand) (results []model_request.Request, err error) {
	log.Println("[service_request.RequestMe] start")
	defer func() {
		log.Println("[service_request.RequestMe] end", "data", map[string]interface{}{"command: ": *c}, "error", err)
	}()

	if err = c.Valid(); err != nil {
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InvalidErr
		service.AddError(ctx, "", "", codeErr)
		return nil, fmt.Errorf(codeErr)
	}

	condition := make(map[string]interface{})
	condition["username"] = c.Username

	cur, err := collection.Request().Collection().Find(ctx, condition)
	if err != nil {
		log.Println("err", err)
		codeErr := src_const.ServiceErr_Request + src_const.ElementErr_Request + src_const.InternalError
		return nil, fmt.Errorf(codeErr)
	}

	results = make([]model_request.Request, 0)
	err = cur.All(ctx, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}
