package cqrs

import "reflect"

type QueryHandler[T, R any] func(T) (R, error)

type QueryRegistry struct {
	queryHandlers map[string]any
}

func NewQueryRegistry() *QueryRegistry {
	return &QueryRegistry{
		queryHandlers: make(map[string]any),
	}
}

func RegisterQueryHandler[T, R any](qr *QueryRegistry, qh QueryHandler[T, R]) {
	var t T
	qr.queryHandlers[reflect.TypeOf(t).String()] = qh
}

func HandleQuery[T, R any](qr *QueryRegistry, q T) (R, error) {
	qh := qr.queryHandlers[reflect.TypeOf(q).String()].(QueryHandler[T, R])
	return qh(q)
}
