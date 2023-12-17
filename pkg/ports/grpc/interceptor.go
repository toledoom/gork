package grpc

import (
	"context"

	"github.com/toledoom/gork/pkg/di"
	"github.com/toledoom/gork/pkg/event"
	"github.com/toledoom/gork/pkg/persistence"
	"google.golang.org/grpc"
)

func WithCommitAndNotifyInterceptor(container *di.Container, setupRepositories persistence.RepositoriesSetup, storageMapper *persistence.StorageMapper) grpc.ServerOption {
	return grpc.UnaryInterceptor(wrapper(container, setupRepositories, storageMapper))
}

func wrapper(container *di.Container, setupRepositories persistence.RepositoriesSetup, storageMapper *persistence.StorageMapper) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		uow := persistence.NewUnitOfWork(storageMapper)
		setupRepositories(container, uow)

		// Calls the handler
		h, err := handler(ctx, req)

		uow.Commit()
		eventPublisher := di.GetService[*event.Publisher](container)

		for _, ev := range uow.DomainEvents() {
			eventPublisher.Notify(ev)
		}

		return h, err
	}
}
