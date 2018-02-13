package command

import "github.com/bwmarrin/discordgo"

// Command ...
type Command interface {
	// Execute will perform the action associated with the command.
	Execute(*Context)
}

// Context ...
type Context struct {
	Prefix  string
	Command string
	Args    []string
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}
