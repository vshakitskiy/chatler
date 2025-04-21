package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "auth.service/api/proto"
	"auth.service/internal/api/handlers"
	"auth.service/internal/repository"
	"auth.service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	grpcServer  *grpc.Server
	port        string
}

func NewApp(
	ctx context.Context,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
) *App {
	port := getEnv("GRPC_PORT", "50051")

	return &App{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		port:        port,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userService := service.NewUserService(a.userRepo)
	authService := service.NewAuthService(
		a.userRepo,
		a.sessionRepo,
		nil,
		time.Duration(0),
		time.Duration(0),
	)
	accessService := service.NewAccessService(authService)

	userHandler := handlers.NewUserServiceHandler(userService)
	authHandler := handlers.NewAuthServiceHandler(authService)
	accessHandler := handlers.NewAccessServiceHandler(accessService)

	a.grpcServer = grpc.NewServer()

	pb.RegisterUserServiceServer(a.grpcServer, userHandler)
	pb.RegisterAuthServiceServer(a.grpcServer, authHandler)
	pb.RegisterAccessServiceServer(a.grpcServer, accessHandler)

	reflection.Register(a.grpcServer)

	go func() {
		lis, err := net.Listen("tcp", ":"+a.port)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		log.Printf("Server gRPC server on port %s", a.port)
		if err := a.grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return a.GracefulShutdown(ctx)
}

func (a *App) GracefulShutdown(ctx context.Context) error {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Shutdown requested via context")
	case <-quit:
		log.Println("Shutdown requested via signal")
	}

	log.Println("Shutting down gRPC server...")
	a.grpcServer.GracefulStop()
	log.Println("gRPC server stopped")

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
