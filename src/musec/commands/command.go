package commands

import (
	"bangarang/musec/instance"
	"bangarang/musec/model"
	"strings"
)

// Commands the bot is capable of performing
const (
	Play = "play"
	Skip = "skip"
	Clear = "clear"
	Leave = "leave"
	Help = "help"
)

// Commands receives parsed text commands from users
var Commands = make(chan *model.Command)

// instances of the bot running on different servers
var instances = make(map[string]*instance.Instance)

func ExecuteCommands() {
	for {
		// Listen on the commands channel
		command := <-Commands

		// Check if we have a session open for this guild
		inst := findOrCreateGuildInstance(command)

		if inst == nil {
			// Skip command as no guild instance was created
			continue
		}

		// TODO change the voice channel if the user is somewhere else

		operand, parameter := parseCommand(command.Message)

		switch operand {
		case Play:
			PlayCommand(inst, parameter)
		case Skip:
			SkipCommand(inst)
		case Clear:
			ClearCommand(inst)
		case Help:
			HelpCommand(inst)
		}
	}
}

// Find or create an instance for the specified guild
func findOrCreateGuildInstance(command *model.Command) *instance.Instance {
	cs := command.Server

	// TODO this should use a mutex

	// Check if we have a session open for this guild
	inst, ok := instances[cs.Guild.ID]

	if !ok {
		// Create a session for the guild
		inst = instance.New(command)

		// If failed to create instance
		if inst == nil {
			return nil
		}

		instances[cs.Guild.ID] = inst
	}

	return inst
}

// parseCommand splits a command into constituent parts
func parseCommand(msg string) (string, string) {
	parts := strings.SplitN(msg, " ", 2)
	operand := strings.ToLower(parts[0])

	// For multipart commands return both parts
	if len(parts) > 1 {
		return operand, parts[1]
	}

	// Else return an operand alone
	return operand, ""
}
