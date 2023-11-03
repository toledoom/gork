package http

import (
	"net/http"

	"github.com/toledoom/gork/pkg/di"
	"github.com/toledoom/gork/pkg/event"
	"github.com/toledoom/gork/pkg/persistence"
)

func WithCommitAndNotifyMiddleware(container *di.Container) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			uow := container.Get("uow")().(persistence.Worker)
			uow.Commit()

			eventPublisher := container.Get("event-publisher")().(*event.Publisher)

			for _, ev := range uow.DomainEvents() {
				eventPublisher.Notify(ev)
			}
		}

		return http.HandlerFunc(fn)
	}
}
