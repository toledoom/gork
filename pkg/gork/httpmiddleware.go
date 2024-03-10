package gork

import (
	"net/http"
)

func WithCommitAndNotifyMiddleware(app *App, storageMapper *StorageMapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			uow := newUnitOfWork(storageMapper)
			app.SetupCommandsAndQueries(uow)

			next.ServeHTTP(w, r)

			uow.Commit()
			eventPublisher := GetService[*EventPublisher](app.container)
			for _, ev := range uow.DomainEvents() {
				eventPublisher.publish(ev)
			}
		}

		return http.HandlerFunc(fn)
	}
}
