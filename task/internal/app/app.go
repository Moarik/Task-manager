package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"taskManager/task/config"
	grpcserver "taskManager/task/internal/adapter/grpc/service"
	httpserver "taskManager/task/internal/adapter/http/service"
	"taskManager/task/internal/adapter/nats/producer"
	"taskManager/task/internal/adapter/postgres"
	"taskManager/task/internal/usecase"
	natsconn "taskManager/task/pkg/nats"
	postgreconn "taskManager/task/pkg/postgre"
)

const serviceName = "task-service"

type App struct {
	httpServer *httpserver.API
	grpcServer *grpcserver.API
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Println(fmt.Sprintf("starting %v service", serviceName))

	log.Println("connecting to postgres", "database", cfg.Postgres.DatabaseName)
	postgresDB, err := postgreconn.NewDB(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	postgresDB.Migrate()

	// nats client
	log.Println("connecting to NATS", "hosts", strings.Join(cfg.Nats.Hosts, ","))
	natsClient, err := natsconn.NewClient(ctx, cfg.Nats.Hosts, cfg.Nats.NKey, cfg.Nats.IsTest)
	if err != nil {
		return nil, fmt.Errorf("nats.NewClient: %w", err)
	}
	log.Println("NATS connection status is", natsClient.Conn.Status().String())

	taskProducer := producer.NewTaskProducer(natsClient, cfg.Nats.NatsSubjects.ClientEventSubject)

	// Repository
	taskRepo := postgres.New(postgresDB)

	// UseCase
	taskUsecase := usecase.NewTask(taskRepo, taskProducer)

	// http service
	httpServer := httpserver.New(cfg.Server, *taskUsecase)

	// grpc service
	gRPCServer := grpcserver.New(
		cfg.Server.GRPCServer,
		taskUsecase,
	)

	app := &App{
		httpServer: httpServer,
		grpcServer: gRPCServer,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.httpServer.Stop()
	if err != nil {
		log.Println("failed to shutdown gRPC service", err)
	}

	err = a.grpcServer.Stop(ctx)
	if err != nil {
		log.Println("failed to shutdown gRPC service", err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	a.httpServer.Run(errCh)
	a.grpcServer.Run(ctx, errCh)

	log.Println(fmt.Sprintf("service %v started", serviceName))

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun

	case s := <-shutdownCh:
		log.Println(fmt.Sprintf("received signal: %v. Running graceful shutdown...", s))

		a.Close(ctx)
		log.Println("graceful shutdown completed!")
	}

	return nil
}
