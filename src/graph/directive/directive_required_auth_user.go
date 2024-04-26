package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nguyencaohiep/metabox-school-proto/golang/authenticator"
	"github.com/nguyencaohiep/metabox-school-proto/grpc_client"

	generated_user "metabox-school-bac-giang-request-service/src/graph/generated/user"
	"metabox-school-bac-giang-request-service/src/network"
)

var UserDirective = generated_user.DirectiveRoot{
	RequiredAuthUser: func(ctx context.Context, obj interface{}, next graphql.Resolver, action *string) (res interface{}, err error) {
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

		if result.Role != RoleStudent {
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
