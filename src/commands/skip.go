package commands

import (
	"musec/discord"
	"musec/instance"
)

func SkipCommand(inst *instance.Instance) {
	sd := inst.State.Server

	// If the skip channel is set
	if inst.Skip != nil {
		discord.TextMessage(sd, "Skipping track... `%d` left in queue", len(inst.Tracks))

		inst.Skip <- true
	}
}
