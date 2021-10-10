package commands

import (
	"bangarang/musec/discord"
	"bangarang/musec/instance"
)

const helpText = `
**Commands**

> play *parameter*

Plays or queues a song or playlist in your current voice channel.
*Note: play can be used with either search terms or a YouTube video URL*

> skip

Skips the current track and begins the next if there is one queued.

> clear

Clears the current queue. 
`

func HelpCommand(inst *instance.Instance)  {
	sd := inst.State.Server

	discord.TextMessage(sd, helpText)
}