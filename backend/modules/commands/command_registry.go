package commands

type CommandHandler func(args string) string

var registry = make(map[string]CommandHandler)

func RegisterCommand(name string, handler CommandHandler) {
	registry[name] = handler
}
