package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"musec/commands"
	"musec/model"
	"musec/service"
	"os"
	"os/signal"
	"syscall"
)

const (
	discordApiEnv = "DISCORD_API_KEY"
	youtubeApiEnv = "YOUTUBE_API_KEY"
)

func main() {
	youtubeApiKey, ytOk := os.LookupEnv(youtubeApiEnv)
	discordApiKey, disOk := os.LookupEnv(discordApiEnv)

	if !ytOk || !disOk {
		log.Printf("You must supply API keys for %s and %s\n", discordApiEnv, youtubeApiEnv)
	}

	service.InitClient(youtubeApiKey)

	log.Println("Created YouTube API client")

	discord, err := discordgo.New("Bot " + discordApiKey)

	if err != nil {
		log.Println("error creating Discord session ", err)
		return
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()

	if err != nil {
		log.Println("error opening WebSocket connection to Discord ", err)
		return
	}

	log.Println("Established connection to Discord")

	// Function which receives commands
	go commands.ExecuteCommands()

	log.Println("Ready for connections")

	// Wait here until CTRL-C or other term signal is received.
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
		log.Println("Error finding guild for id " + msg.GuildID)
	}

	textChannel, err := s.Channel(msg.ChannelID)

	if err != nil {
		log.Println("Error finding channel for id " + msg.ChannelID)
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
		log.Println("Error closing session ", err)
	}

	os.Exit(0)
}
