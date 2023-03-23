package commands

import (
	"bangarang/musec/discord"
	"bangarang/musec/instance"
	"bangarang/musec/model"
	"bangarang/musec/service"
	"fmt"
	"net/url"
	"regexp"
)

const listGroup = "list"

var playlistRegex, _ = regexp.Compile("list=(?P<list>[a-zA-Z0-9_-]+)")

func PlayCommand(inst *instance.Instance, param string) {
	sd := inst.State.Server

	// Return early if no parameter was supplied
	if param == "" {
		discord.TextMessage(sd, "Use like `play rick astley` or `play https://www.youtube.com/watch?v=dQw4w9WgXcQ`")

		return
	}

	if isUrl(param) {
		// First attempt to extract a playlist id
		matches := playlistRegex.FindStringSubmatch(param)

		index := playlistRegex.SubexpIndex(listGroup)

		// If found then try and play the playlist
		if len(matches) > index && index != -1 {
			playlistId := matches[index]

			// Retrieve the videoIds
			ids, err := service.VideosFromPlaylist(playlistId)

			if err != nil {
				// Try and play as a video
				err := playTrack(inst, param)

				if err != nil {
					// Only complain if that fails as well
					discord.TextMessage(sd, fmt.Sprintf("Nothing found for playlist `%s`", param))
				}

				return
			}

			playPlaylist(inst, ids)

			return
		}

		// Otherwise, try and play the single url
		_ = playTrack(inst, param)
	} else {
		// Search using the parameter as search terms
		result, err := service.SearchYouTube(param)

		if err != nil {
			discord.TextMessage(sd, fmt.Sprintf("No results found for search `%s`", param))

			// Return without queueing the track
			return
		}

		if result.Kind == service.KindPlaylist {
			playPlaylist(inst, result)
		} else {
			_ = playTrack(inst, result.Id)
		}
	}
}

func playPlaylist(inst *instance.Instance, playlist *service.SearchResult) {
	sd := inst.State.Server
	ch := make(chan bool)

	for _, vr := range playlist.Videos {
		vr := vr

		go func() {
			video, format, err := service.GetYoutubeVideo(vr.Id)

			// Increment count if no error
			if err != nil {
				println("failed to add a track from a playlist " + playlist.Title)
			} else {
				stop := make(chan bool)

				track := &model.QueuedTrack{
					Video:  video,
					Format: format,
					Stop:   stop,
				}

				// Queue the track up
				inst.Tracks <- track
			}
			ch <- true
		}()
	}

	for _ = range playlist.Videos {
		// Discard channel items
		<-ch
	}

	size := len(inst.Tracks)

	discord.TextMessage(sd, fmt.Sprintf("Added `%d` tracks from `%s`, total `%d` in the queue", len(playlist.Videos), playlist.Title, size))
}

func playTrack(inst *instance.Instance, id string) error {
	sd := inst.State.Server
	video, format, err := service.GetYoutubeVideo(id)

	if err != nil {
		discord.TextMessage(sd, fmt.Sprintf("No video found at id/URL `%s`", id))

		// Return without queueing the track
		return err
	}

	stop := make(chan bool)

	track := &model.QueuedTrack{
		Video:  video,
		Format: format,
		Stop:   stop,
	}

	pos := len(inst.Tracks) + 1

	// Only send queued message if not the first item
	if pos > 1 || inst.Skip != nil {
		discord.TextMessage(sd, fmt.Sprintf("Adding `%s` as number `%d` in the queue", video.Title, pos))
	}

	// Queue the track up
	inst.Tracks <- track

	return nil
}

// Roughly checks if the parameter is a URL
func isUrl(param string) bool {
	u, err := url.Parse(param)

	if err != nil {
		return false
	}

	return u.Host != "" && u.Path != ""
}
