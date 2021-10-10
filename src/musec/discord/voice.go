package discord

import (
	"bangarang/musec/model"
	"github.com/bwmarrin/discordgo"
)

func JoinVoiceChannel(state *model.ConnectionState) (*discordgo.VoiceConnection, error) {
	sd := state.Server
	s, vc := sd.Session, state.VoiceChannel

	// Join the voice channel muted
	conn, err := s.ChannelVoiceJoin(vc.GuildID, vc.ID, false, true)

	return conn, err
}

func FindUserVoiceChannel(sd *model.ServerDetails, u *discordgo.User) *discordgo.Channel {
	s, g := sd.Session, sd.Guild

	// Loop through active voice connections to find user
	for _, vs := range g.VoiceStates {
		if u.ID == vs.UserID {
			// Find channel in server
			vc, err := s.Channel(vs.ChannelID)

			if err != nil {
				println("Error finding channel for user " + u.Username)

				return nil
			}

			return vc
		}
	}

	return nil
}
