package commands

import (
	"bangarang/musec/discord"
	"bangarang/musec/instance"
	"fmt"
)

func ClearCommand(inst *instance.Instance) {
	size := len(inst.Tracks)

	discord.TextMessage(inst.State.Server, fmt.Sprintf("Clearing `%d` tracks from the queue", size))

	// Read messages till queue is empty
	for len(inst.Tracks) > 0 {
		<-inst.Tracks
	}
}
