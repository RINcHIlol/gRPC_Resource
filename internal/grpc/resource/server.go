package resource

import (
	"context"
	"fmt"
	resourceService "github.com/RINcHIlol/protosFirst/gen/go/resourceService"
	"google.golang.org/grpc"
)

type ResourceServer interface {
	AccessResource(ctx context.Context, resource_name string) (bool, string, error)
	LogAccessAttempt(ctx context.Context, user_id int64, resource_name string, access bool) (string, error)
}

type serverApi struct {
	//запуск приложения без реализации всех методов интерфейса
	resourceService.UnimplementedResourceServer
	resourceServer ResourceServer
}

func Register(gRPC *grpc.Server, resourceServer ResourceServer) {
	resourceService.RegisterResourceServer(gRPC, &serverApi{resourceServer: resourceServer})
}

func (s *serverApi) AccessResource(ctx context.Context, request *resourceService.AccessResourceRequest) (*resourceService.AccessResourceResponse, error) {
	const op = "resource.AccessResource"
	fmt.Printf("%s: %s", op, "accessResource\n")

	access, message, err := s.resourceServer.AccessResource(ctx, request.ResourceName)
	if err != nil {
		fmt.Printf("%s, %s", op, err)
		return nil, err
	}

	return &resourceService.AccessResourceResponse{Access: access, Message: message}, nil
}

func (s *serverApi) LogAccessAttempt(ctx context.Context, request *resourceService.LogAccessAttemptRequest) (*resourceService.LogAccessAttemptResponse, error) {
	const op = "resource.LogAccessAttempt"
	fmt.Printf("%s: %s", op, "LogAccessAttempt")

	return &resourceService.LogAccessAttemptResponse{}, nil
}
