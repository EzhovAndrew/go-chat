package proxy

import (
	"context"
	"log"

	"github.com/go-chat/gateway/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

func RegisterServices(ctx context.Context, mux *runtime.ServeMux, cfg *config.Config) error {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.Services.Auth, opts); err != nil {
		return err
	}
	log.Println("Registered Auth Service")

	if err := usersv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.Services.Users, opts); err != nil {
		return err
	}
	log.Println("Registered Users Service")

	if err := chatv1.RegisterChatServiceHandlerFromEndpoint(ctx, mux, cfg.Services.Chat, opts); err != nil {
		return err
	}
	log.Println("Registered Chat Service")

	if err := socialv1.RegisterSocialServiceHandlerFromEndpoint(ctx, mux, cfg.Services.Social, opts); err != nil {
		return err
	}
	log.Println("Registered Social Service")

	if err := notificationsv1.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, cfg.Services.Notifications, opts); err != nil {
		return err
	}
	log.Println("Registered Notifications Service")

	return nil
}

