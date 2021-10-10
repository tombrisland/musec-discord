package commands

import (
	"bangarang/musec/discord"
	"bangarang/musec/instance"
	"fmt"
)

func SkipCommand(inst *instance.Instance) {
	sd := inst.State.Server

	// If the skip channel is set
	if inst.Skip != nil {
		discord.TextMessage(sd, fmt.Sprintf("Skipping track... `%d` left in queue", len(inst.Tracks)))

		inst.Skip <- true
	}
}
