package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/places-svc/internal/api/grpc/meta"
	"github.com/chains-lab/places-svc/internal/config/constant"
	"github.com/chains-lab/places-svc/internal/problems"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// EXAMPLE
var userAllowList = map[string]map[string]bool{
	//// Employee methods
	//employee.EmployeeService_GetDistributorEmployee_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//employee.EmployeeService_SelectEmployees_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//employee.EmployeeService_GetInvite_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//employee.EmployeeService_SelectInvites_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//
	//employee.EmployeeService_UpdateEmployeeRole_FullMethodName: {roles.User: true},
	//employee.EmployeeService_DeleteEmployee_FullMethodName:     {roles.User: true},
	//employee.EmployeeService_SendInvite_FullMethodName:         {roles.User: true},
	//employee.EmployeeService_WithdrawInvite_FullMethodName:     {roles.User: true},
	//employee.EmployeeService_AcceptInvite_FullMethodName:       {roles.User: true},
	//employee.EmployeeService_RejectInvite_FullMethodName:       {roles.User: true},
	//
	//// Distributor methods
	//distributor.DistributorService_GetDistributorBlock_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//distributor.DistributorService_SelectDistributorBlocks_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//distributor.DistributorService_GetDistributor_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//distributor.DistributorService_SelectDistributors_FullMethodName: {
	//	roles.SuperUser: true,
	//	roles.Admin:     true,
	//	roles.Moder:     true,
	//	roles.User:      true,
	//},
	//
	//distributor.DistributorService_CreateDistributor_FullMethodName:            {roles.User: true},
	//distributor.DistributorService_UpdateDistributorName_FullMethodName:        {roles.User: true},
	//distributor.DistributorService_UpdateDistributorIcon_FullMethodName:        {roles.User: true},
	//distributor.DistributorService_SetDistributorStatusInactive_FullMethodName: {roles.User: true},
	//distributor.DistributorService_SetDistributorStatusActive_FullMethodName:   {roles.User: true},
	//
	////Admin methods
	//distributor.DistributorService_BlockDistributor_FullMethodName:   {roles.SuperUser: true, roles.Admin: true},
	//distributor.DistributorService_UnblockDistributor_FullMethodName: {roles.SuperUser: true, roles.Admin: true},
}

// EXAMPLE
var serviceAllowList = map[string]bool{
	//employee.EmployeeService_GetDistributorEmployee_FullMethodName: true,
	//employee.EmployeeService_SelectEmployees_FullMethodName:        true,
	//employee.EmployeeService_GetInvite_FullMethodName:              true,
	//employee.EmployeeService_SelectInvites_FullMethodName:          true,
	//
	//distributor.DistributorService_GetDistributorBlock_FullMethodName:     true,
	//distributor.DistributorService_SelectDistributorBlocks_FullMethodName: true,
	//distributor.DistributorService_GetDistributor_FullMethodName:          true,
	//distributor.DistributorService_SelectDistributors_FullMethodName:      true,
}

func Auth(skUser, skService string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("no metadata found in incoming context"))
		}

		srvToken := md.Get("x-service-token")
		if len(srvToken) == 0 {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("service token not supplied"))
		}

		data, err := auth.VerifyServiceJWT(ctx, srvToken[0], skService)
		if err != nil {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("failed to verify service token"))
		}

		thisSvcInAudience := false
		for _, aud := range data.Audience {
			if aud == constant.ServiceName {
				thisSvcInAudience = true
				break
			}
		}

		if !thisSvcInAudience {
			return nil, status.New(codes.Unauthenticated,
				fmt.Sprintf("service issuer %s not in audience %v", data.Issuer, data.Audience)).Err()
		}

		userToken := md.Get("x-user-token")

		if serviceAllowList[info.FullMethod] && len(userToken) > 0 {
			if _, alsoUser := userAllowList[info.FullMethod]; !alsoUser {
				return nil, problems.RaisePermissionDenied(ctx, fmt.Errorf("method is service-only"))
			}
		}

		if len(userToken) == 0 {
			if serviceAllowList[info.FullMethod] {
				return handler(ctx, req)
			}

			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("user token not supplied"))
		}

		userData, err := auth.VerifyUserJWT(ctx, userToken[0], skUser)
		if err != nil {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("failed to verify user token"))
		}

		userID, err := uuid.Parse(userData.Subject)
		if err != nil {
			return nil, problems.RaiseUnauthenticated(ctx, fmt.Errorf("invalid user id: %v", err))
		}

		ctx = context.WithValue(ctx, meta.UserCtxKey, meta.UserData{
			ID:        userID,
			SessionID: userData.Session,
			Verified:  userData.Verified,
			Role:      userData.Role,
		})

		if rolesForMethod, ok := userAllowList[info.FullMethod]; ok {
			user, err := meta.User(ctx)
			if err != nil {
				return nil, err
			}
			if !rolesForMethod[user.Role] {
				return nil, problems.RaisePermissionDenied(ctx, fmt.Errorf("user does not have required permissions for this method"))
			}
		}

		return handler(ctx, req)
	}
}
