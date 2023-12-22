package cqrs

import (
	"fmt"
	"reflect"
)

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

type QueryNotRegisteredError struct {
	q interface{}
}

func (qnre *QueryNotRegisteredError) Error() string {
	return fmt.Sprintf("query handler not registered for query %s", reflect.TypeOf(qnre.q).String())
}

func HandleQuery[T, R any](qr *QueryRegistry, q T) (R, error) {
	tryQueryHandler, ok := qr.queryHandlers[reflect.TypeOf(q).String()]
	if !ok {
		var r R
		return r, &QueryNotRegisteredError{q: q}
	}

	qh := tryQueryHandler.(QueryHandler[T, R])
	return qh(q)
}
