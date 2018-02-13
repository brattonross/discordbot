package command

// Handler handles a number of commands.
type Handler struct {
	commands map[string]Command
}

// NewHandler creates a new instance of Handler.
func NewHandler() *Handler {
	return &Handler{
		commands: make(map[string]Command),
	}
}

// AddCommand adds the given command to the handler's commands.
func (h *Handler) AddCommand(name string, c Command) {
	h.commands[name] = c
}

// RemoveCommand removes the command with the given name.
func (h *Handler) RemoveCommand(name string) {
	delete(h.commands, name)
}

// Command returns the command with the given name.
// If the command does not exist, nil is returned instead.
func (h Handler) Command(name string) Command {
	return h.commands[name]
}
