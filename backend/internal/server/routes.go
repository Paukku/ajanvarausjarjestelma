package server

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	pbcommon "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
	pbHTTP "github.com/Paukku/ajanvarausjarjestelma/backend/pb/http"
)

type ApiRegister struct {
	rule func(
		cb func(ctx context.Context, w http.ResponseWriter, r *http.Request,
			arg proto.Message, ret proto.Message, err error),
		interceptors ...grpc.UnaryServerInterceptor,
	) (method string, path string, handler http.HandlerFunc)
	accessRole pbcommon.UserRole
}

func RegisterRoutes(mux *http.ServeMux, userConverter *pbHTTP.BusinessCustomerAPIHTTPConverter) {
	userApiRegister := []ApiRegister{
		{rule: userConverter.CreateUserHTTPRule, accessRole: pbcommon.UserRole_OWNER},
		{rule: userConverter.GetUserHTTPRule, accessRole: pbcommon.UserRole_OWNER},
		{rule: userConverter.GetUserByIdHTTPRule, accessRole: pbcommon.UserRole_OWNER},
	}

	// TULEVA ADMIN LISTA
	// adminApiRegister := []ApiRegister{
	// 	  {rule: adminConverter.CreateCompanyHTTPRule, accessRole: pbcommon.UserRole_ADMIN},
	// }

	for _, api := range userApiRegister {
		_, path, handler := api.rule(nil)

		fmt.Println("Registering route:", path)

		// wrapataan middlewareen
		mux.Handle(path, RoleMiddleware(api.accessRole)(handler))
	}

}
