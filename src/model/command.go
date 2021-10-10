package model

import "github.com/bwmarrin/discordgo"

// Command to be parsed as a command
type Command struct {
	// Server state of the connection
	Server *ServerDetails

	// User who sent the message
	User *discordgo.User

	// Message content
	Message string
}
