package service

import (
	"bangarang/musec/model"
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"io"
)

// Mime type matching audio streams
const audioMime = "audio"

// The dlClient to talk to YouTube
var dlClient = youtube.Client{}

// The downloader to fetch videos
var download = downloader.Downloader{
	Client: dlClient,
	// Intercepted before it writes to disk
	OutputDir: "",
}

// GetYoutubeVideo searches the URL and returns service metadata
func GetYoutubeVideo(url string) (*youtube.Video, *youtube.Format, error) {
	println("searching for video with url " + url)

	video, err := dlClient.GetVideo(url)

	if err != nil {
		return nil, nil, err
	}

	println("found video with title " + video.Title)

	// Select the best audio stream
	formats := video.Formats.Type(audioMime)
	formats.Sort()

	if len(formats) == 0 {
		return nil, nil, errors.New("no audio formats found")
	}

	format := formats[0]

	fmt.Printf("chose format %v\n", format.AudioQuality)

	return video, &format, nil
}

// GetYoutubeAudioStream returns an io.ReadCloser representing the audio stream
func GetYoutubeAudioStream(track *model.QueuedTrack) (io.ReadCloser, error) {
	stream, _, err := download.GetStream(track.Video, track.Format)

	if err != nil {
		println("failed to get audio stream for " + track.Video.Title, err)

		return nil, errors.New("failed to get audio stream")
	}

	return stream, err
}
