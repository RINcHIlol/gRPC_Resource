package resourceService

import (
	"ResourceService/internal/domain/models"
	"context"
	"fmt"
	authpb "github.com/RINcHIlol/protosFirst/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"log/slog"
	"time"
)

type ResourceService struct {
	log                *slog.Logger
	zaglushkaInterface DataProvider
	tokenTTL           time.Duration
}

type DataProvider interface {
	Data(ctx context.Context, resourceName string) (models.Data, error)
}

func New(log *slog.Logger, zaglushkaInterface DataProvider, tokenTTL time.Duration) *ResourceService {
	return &ResourceService{
		log:                log,
		zaglushkaInterface: zaglushkaInterface,
		tokenTTL:           tokenTTL,
	}
}

func (r *ResourceService) AccessResource(ctx context.Context, resourceName string) (bool, string, error) {
	const op = "resourceService.AccessResource"
	r.log.Info("%s: %s", op, "accessResource")

	conn, err := grpc.Dial("localhost:44044", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения к AuthService: %v", err)
	}
	defer conn.Close()
	authClient := authpb.NewAuthClient(conn)

	loginResp, err := authClient.Login(ctx, &authpb.LoginRequest{
		Email:    "user@example.com",
		Password: "password",
		AppId:    1,
	})
	if err != nil {
		log.Fatalf("Ошибка при авторизации: %v", err)
	}
	jwtToken := loginResp.Token

	md := metadata.New(map[string]string{"authorization": "Bearer " + jwtToken})
	meta := metadata.NewOutgoingContext(ctx, md)

	// change id to nonAdmin user
	adminResp, err := authClient.IsAdmin(meta, &authpb.IsAdminRequest{UserId: 6})
	if err != nil {
		log.Fatalf("Ошибка при проверке прав администратора: %v", err)
	}

	fmt.Println("Пользователь админ?", adminResp.IsAdmin)

	data, err := r.zaglushkaInterface.Data(ctx, resourceName)
	if err != nil {
		r.log.Error(op, err)
		return false, "", err
	}

	if data.Success == false && adminResp.IsAdmin == false {
		data.Info = "недосутпно"
	} else {
		data.Success = true
	}

	return data.Success, data.Info, nil
}

func (r *ResourceService) LogAccessAttempt(ctx context.Context, userId int64, resourceName string, access bool) (string, error) {
	const op = "resourceService.LogAccessAttempt"
	r.log.Info("%s: %s", op, "logAccessAttempt")

	return "", nil
}
