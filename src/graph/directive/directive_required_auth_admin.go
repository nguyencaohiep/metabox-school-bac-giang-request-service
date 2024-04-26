package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nguyencaohiep/metabox-school-proto/golang/authenticator"
	"github.com/nguyencaohiep/metabox-school-proto/grpc_client"

	generatedAdmin "metabox-school-bac-giang-request-service/src/graph/generated/admin"
	"metabox-school-bac-giang-request-service/src/network"
)

const (
	RoleSuperAdmin = "super-admin"
	RoleAdmin      = "admin"
	RoleStudent    = "student"
)

var AdminDirective = generatedAdmin.DirectiveRoot{
	RequiredAuthAdmin: func(ctx context.Context, obj interface{}, next graphql.Resolver, action []*string, actionAdmin *string, check_ip *bool) (res interface{}, err error) {
		if !network.HasToken(ctx) {
			return nil, fmt.Errorf("unauthorized")
		}

		tokenStr := network.Token(ctx)
		// result, err := service_shared.TokenVerify(ctx, &authenticator.TokenVerifyRequest{JwtToken: tokenStr})
		result, err := grpc_client.AuthenticatorClient().TokenVerify(ctx, &authenticator.TokenVerifyRequest{
			JwtToken: tokenStr,
		})
		if err != nil || result == nil {
			return nil, err
		}

		if result.Role != RoleSuperAdmin {
			return nil, fmt.Errorf("permission deny")
		}

		ctx = context.WithValue(ctx, "account_id", result.AccountId)
		ctx = context.WithValue(ctx, "role", result.Role)
		ctx = context.WithValue(ctx, "username", result.UserName)

		if action == nil {
			return next(ctx)
		}

		return next(ctx)
	},
	RequiredAuthSuperUser: func(ctx context.Context, obj interface{}, next graphql.Resolver, action *string) (res interface{}, err error) {
		if !network.HasToken(ctx) {
			return nil, fmt.Errorf("unauthorized")
		}

		tokenStr := network.Token(ctx)
		// result, err := service_shared.TokenVerify(ctx, &authenticator.TokenVerifyRequest{JwtToken: tokenStr})
		result, err := grpc_client.AuthenticatorClient().TokenVerify(ctx, &authenticator.TokenVerifyRequest{
			JwtToken: tokenStr,
		})
		if err != nil || result == nil {
			return nil, err
		}

		if result.Role != RoleSuperAdmin && result.Role != RoleAdmin {
			return nil, fmt.Errorf("permission deny")
		}

		ctx = context.WithValue(ctx, "account_id", result.AccountId)
		ctx = context.WithValue(ctx, "role", result.Role)
		ctx = context.WithValue(ctx, "username", result.UserName)

		if action == nil {
			return next(ctx)
		}

		return next(ctx)
	},
}
