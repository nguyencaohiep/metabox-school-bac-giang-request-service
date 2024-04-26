package resolver_admin

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	generated_admin "metabox-school-bac-giang-request-service/src/graph/generated/admin"
	graph_model "metabox-school-bac-giang-request-service/src/graph/generated/model"
	service_request "metabox-school-bac-giang-request-service/src/service/request"

	"github.com/99designs/gqlgen/graphql"
)

// RequestPagination is the resolver for the requestPagination field.
func (r *queryResolver) RequestPagination(ctx context.Context, page int, limit int, orderBy *string, search map[string]interface{}) (*graph_model.RequestPagination, error) {
	input := &service_request.RequestPaginationCommand{
		Page:  page,
		Limit: limit,
	}

	if orderBy != nil {
		input.OrderBy = *orderBy
	}

	if search != nil {
		input.Search = search
	}

	total, result, err := service_request.RequestPagination(ctx, input)
	if err != nil {
		if graphql.GetErrors(ctx) == nil {
			return nil, err
		}
		return nil, nil
	}

	requests := make([]graph_model.Request, 0)
	for i := 0; i < len(result); i++ {
		requests = append(requests, *result[i].ConvertToModelGraph())
	}

	return &graph_model.RequestPagination{
		Rows: requests,
		Paging: graph_model.Pagination{
			CurrentPage: page,
			Limit:       limit,
			TotalPage:   CalculateTotalPage(total, limit),
			Total:       total,
		},
	}, nil
}

// Query returns generated_admin.QueryResolver implementation.
func (r *Resolver) Query() generated_admin.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
