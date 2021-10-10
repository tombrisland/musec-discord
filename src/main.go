package main

import (
	"bangarang/musec/commands"
	"bangarang/musec/model"
	"bangarang/musec/service"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

const (
	discordApiEnv = "DISCORD_API_KEY"
	youtubeApiEnv = "YOUTUBE_API_KEY"
)

func main() {
	youtubeApiKey, b := os.LookupEnv(youtubeApiEnv)
	discordApiKey, b := os.LookupEnv(discordApiEnv)

	if !b {
		println(fmt.Sprintf("You must supply API keys for %s and %s", discordApiEnv, youtubeApiEnv))
	}

	service.InitClient(youtubeApiKey)

	discord, err := discordgo.New("Bot " + discordApiKey)

	if err != nil {
		fmt.Println("error creating Discord session ", err)
		return
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()

	if err != nil {
		fmt.Println("error opening WebSocket connection to Discord ", err)
		return
	}

	// Function which receives commands
	go commands.ExecuteCommands()

	// Wait here until CTRL-C or other term signal is received.
	println("bot is now running...")

	endOnInterrupt(discord)
}

// messageCreate handler when a message is received
func messageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	author := msg.Author

	// Ignore messages sent by bot
	if author.ID == s.State.User.ID {
		return
	}

	// Retrieve the guild from state
	guild, err := s.State.Guild(msg.GuildID)

	if err != nil {
		println("Error finding guild for id " + msg.GuildID)
	}

	textChannel, err := s.Channel(msg.ChannelID)

	if err != nil {
		println("Error finding channel for id " + msg.ChannelID)
	}

	// Create a new command with the message details
	cs := &model.ServerDetails{
		Session:     s,
		Guild:       guild,
		TextChannel: textChannel,
	}

	command := &model.Command{
		Server:  cs,
		User:    author,
		Message: msg.Content,
	}

	// Send the command for execution
	commands.Commands <- command
}

// Run until an interrupt is received
func endOnInterrupt(session *discordgo.Session) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err := session.Close()

	if err != nil {
		println("Error closing session ", err)
	}

	os.Exit(0)
}
