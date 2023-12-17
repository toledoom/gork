package gork

import (
	"net/http"

	"github.com/toledoom/gork/pkg/event"
)

func WithCommitAndNotifyMiddleware(container *Container, setupRepositories RepositoriesSetup, storageMapper *StorageMapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			uow := NewUnitOfWork(storageMapper)
			setupRepositories(container, uow)

			next.ServeHTTP(w, r)

			uow.Commit()
			eventPublisher := GetService[*event.Publisher](container)
			for _, ev := range uow.DomainEvents() {
				eventPublisher.Notify(ev)
			}
		}

		return http.HandlerFunc(fn)
	}
}
