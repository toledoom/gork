package cqrs

type Query interface {
	QueryID() string
}

type QueryHandler interface {
	Handle(c Query) (QueryResponse[any], error)
	QueryID() string
}

type QueryResponse[T any] interface {
	Data() T
}

type QueryBus struct {
	queries map[string]QueryHandler
}

func NewQueryBus(queryHandlerList []QueryHandler) *QueryBus {
	queryHandlerMap := make(map[string]QueryHandler)
	for _, qh := range queryHandlerList {
		queryHandlerMap[qh.QueryID()] = qh
	}
	return &QueryBus{
		queries: queryHandlerMap,
	}
}

func (qb *QueryBus) Handle(q Query) (QueryResponse[any], error) {
	qh := qb.queries[q.QueryID()]
	return qh.Handle(q)
}
