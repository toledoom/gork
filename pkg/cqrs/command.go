package cqrs

type Command interface {
	CmdID() string
}

type CommandHandler interface {
	Handle(c Command) error
	CmdID() string
}

type CommandBus struct {
	commands map[string]CommandHandler
}

func NewCommandBus(commandHandlerList []CommandHandler) *CommandBus {
	commandMap := make(map[string]CommandHandler)
	for _, ch := range commandHandlerList {
		commandMap[ch.CmdID()] = ch
	}
	return &CommandBus{
		commands: commandMap,
	}
}

func (cb *CommandBus) Handle(c Command) error {
	ch := cb.commands[c.CmdID()]
	return ch.Handle(c)
}
