package resource

import (
	"context"
	resourceServicev1 "github.com/RINcHIlol/protosFirst/gen/go/resourceService"
)

type ResourceServer interface {
	AccessResource(ctx context.Context, resource_name string) (bool, string, error)
	LogAccessAttempt(ctx context.Context, user_id int64, resource_name string, access bool) (string, error)
}

type serverApi struct {
	//запуск приложения без реализации всех методов интерфейса
	resourceServicev1.UnimplementedResourceServer
	resourceServer ResourceServer
}
