package gork

import (
	"context"

	"google.golang.org/grpc"
)

func WithCommitAndNotifyInterceptor(app *App, storageMapper *StorageMapper) grpc.ServerOption {
	return grpc.UnaryInterceptor(wrapper(app, storageMapper))
}

func wrapper(app *App, storageMapper *StorageMapper) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		uow := newUnitOfWork(storageMapper)
		app.SetupCommandsAndQueries(uow)

		// Calls the handler
		h, err := handler(ctx, req)

		uow.Commit()
		eventPublisher := GetService[*EventPublisher](app.container)

		for _, ev := range uow.DomainEvents() {
			eventPublisher.publish(ev)
		}

		return h, err
	}
}
