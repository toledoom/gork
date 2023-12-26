package gork_test

type dumbCommand struct {
	ID string
}

func dumbCommandHandler(dc *dumbCommand) error { return nil }

type dumbQuery struct{}

func dumbQueryHandler(dc *dumbQuery) (string, error) { return "a value", nil }
