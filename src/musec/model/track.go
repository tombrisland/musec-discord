package model

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

// QueuedTrack representing a single track for instance
type QueuedTrack struct {
	// Video the YouTube service metadata
	Video *youtube.Video

	// Stop channel to cease instance of this track
	Stop chan bool

	// Format of the audio stream
	Format *youtube.Format

	// User who requested the track
	User *discordgo.User
}
