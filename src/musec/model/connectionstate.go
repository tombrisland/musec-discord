package model

import "github.com/bwmarrin/discordgo"

// ServerDetails containing connection information
type ServerDetails struct {
	// Session for the discord API
	Session *discordgo.Session

	// Guild we are currently connected to
	Guild *discordgo.Guild

	// TextChannel to speak in
	TextChannel *discordgo.Channel
}

// ConnectionState of the voice connection to a server
type ConnectionState struct {
	// Server state of the connection
	Server *ServerDetails

	// VoiceChannel to play the tracks in
	VoiceChannel *discordgo.Channel

	// VoiceConnection if we are in the voice channel already
	VoiceConnection *discordgo.VoiceConnection
}
