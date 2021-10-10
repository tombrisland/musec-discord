package commands

import (
	"musec/discord"
	"musec/instance"
)

func ClearCommand(inst *instance.Instance) {
	size := len(inst.Tracks)

	discord.TextMessage(inst.State.Server, "Clearing `%d` tracks from the queue", size)

	// Read messages till queue is empty
	for len(inst.Tracks) > 0 {
		<-inst.Tracks
	}
}
