package gork

import (
	"context"

	"google.golang.org/grpc"
)

func WithCommitAndNotifyInterceptor(container *Container, setupRepositories RepositoriesSetup, storageMapper *StorageMapper) grpc.ServerOption {
	return grpc.UnaryInterceptor(wrapper(container, setupRepositories, storageMapper))
}

func wrapper(container *Container, setupRepositories RepositoriesSetup, storageMapper *StorageMapper) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		uow := NewUnitOfWork(storageMapper)
		setupRepositories(container, uow)

		// Calls the handler
		h, err := handler(ctx, req)

		uow.Commit()
		eventPublisher := GetService[*EventPublisher](container)

		for _, ev := range uow.DomainEvents() {
			eventPublisher.Publish(ev)
		}

		return h, err
	}
}
