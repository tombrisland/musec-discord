package instance

import (
	"bangarang/musec/discord"
	"bangarang/musec/model"
	"bangarang/musec/service"
	"fmt"
	"time"
)

// maxQueueSize of the tracks for instance
const maxQueueSize = 10_000

// afkDelay after which the bot leaves the channel
const afkDelay = 5 * time.Minute

type Instance struct {
	// Tracks to be played
	Tracks chan *model.QueuedTrack

	// Skip message to skip the current track
	Skip chan bool

	// State of the current connection to a discord server
	State *model.ConnectionState
}


// New Instance for a specific guild connection
func New(command *model.Command) *Instance {
	// Get the channel the user is currently in
	vc := discord.FindUserVoiceChannel(command.Server, command.User)

	// State object for this guild
	connState := &model.ConnectionState{
		Server: command.Server,
		VoiceChannel: vc,
		// No voice connection yet
		VoiceConnection: nil,
	}

	// Channel for queueing tracks
	tracks := make(chan *model.QueuedTrack, maxQueueSize)

	playback := Instance{Tracks: tracks, Skip: nil, State: connState}

	// Start the instance function
	go playback.BeginPlayback()

	return &playback
}

// BeginPlayback listens to the tracks channel and plays them as they come in
func (p *Instance) BeginPlayback() {
	sd := p.State.Server
	cs := p.State

	// Loop until channel is closed
	for {
		// Leave the channel on sustained inactivity
		if cs.VoiceConnection != nil && p.Inactive() {
			go p.leave()
		}

		track := <-p.Tracks

		// Stop this specific track if skip is called
		p.Skip = track.Stop

		if cs.VoiceConnection == nil || !cs.VoiceConnection.Ready {
			// If there is no voice channel currently set
			if cs.VoiceChannel == nil {
				// Try and join the user who requested the track
				vc := discord.FindUserVoiceChannel(sd, track.User)

				// If they aren't in a channel
				if vc == nil {
					discord.TextMessage(sd, "Get in a voice channel " + track.User.Mention())

					// Skip and wait for next track to try again
					continue
				}

				cs.VoiceChannel = vc
			}
			vc, err := discord.JoinVoiceChannel(cs)

			// Set the state to have the active voice connection
			cs.VoiceConnection = vc

			if err != nil {
				discord.TextMessage(sd, "I couldn't join " + cs.VoiceChannel.Mention())

				// Skip and wait for next track to try again
				continue
			}
		}

		// Get the audio stream for that track
		stream, err := service.GetYoutubeAudioStream(track)

		if err != nil {
			discord.TextMessage(sd, fmt.Sprintf("I wasn't able to play `%s`", track.Video.Title))

			continue
		}

		discord.TextMessage(sd, fmt.Sprintf("Playing `%s`", track.Video.Title))

		// Play the stream until it ends or a stop message is received
		err = PlayAudioStream(cs.VoiceConnection, stream, track.Stop)

		if err != nil {
			discord.TextMessage(sd, fmt.Sprintf("I had trouble while playing `%s`", track.Video.Title))
		}

		// Nil the skip channel for now
		p.Skip = nil
	}
}

func (p *Instance) Inactive() bool {
	return len(p.Tracks) == 0 && p.Skip == nil
}

// leave the voice channel after sustained inactivity
func (p *Instance) leave() {
	// Sleep for a while
	time.Sleep(afkDelay)

	s := p.State

	if p.Inactive() {
		println("leaving voice channel")

		// Leave the channel until re-activated
		err := s.VoiceConnection.Disconnect()

		discord.TextMessage(s.Server, fmt.Sprintf("Left voice after `%ds` of inactivity", int(afkDelay.Seconds())))

		if err != nil {
			println("failed to leave voice channel")
		}
	}
}