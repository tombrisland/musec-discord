package service

import (
	"context"
	"errors"
	"google.golang.org/api/option"
)

import "google.golang.org/api/youtube/v3"

// Return snippetType in api response
const snippetType = "snippet"

var snippetParam = []string{snippetType}

const (
	KindVideo    = "youtube#video"
	KindPlaylist = "youtube#playlist"
)

// Count of maxSearchResults to return in a search
const maxSearchResults = 1

// Count of maxPlaylistResults to return from a playlist
const maxPlaylistResults = 50

// The firstResult of the search
const firstResult = 0

// client for calling YouTube search and playlists
var client *youtube.Service = nil

func InitClient(apiKey string) {
	ctx := context.Background()
	// Try to initialise with YouTube API key
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		println("Failed to initialise YouTube search client", err.Error())
	}

	client = service
}

type SearchResult struct {
	// Title of the video or playlist
	Title string

	// Id for the video or playlist
	Id string

	// Kind either video or playlist
	Kind string

	// Videos if this is a playlist
	Videos []SearchResult
}

// SearchYouTube and return a videoId or videoIds in the case of a playlist
func SearchYouTube(terms string) (*SearchResult, error) {
	req := client.Search.List([]string{snippetType}).Q(terms).MaxResults(maxSearchResults)

	res, err := req.Do()

	if err != nil {
		println("failed to search YouTube for terms " + terms)

		return nil, err
	}

	items := res.Items

	if len(items) < 1 {
		println("no result found for terms " + terms)

		return nil, errors.New("couldn't find a result for " + terms)
	}

	// Select the first search result
	item := items[firstResult]

	// Set the title and kind as they are common
	result := SearchResult{
		Title: item.Snippet.Title,
		Kind:  item.Kind,
	}

	if item.Id.Kind == KindVideo {
		// Set the videoId
		result.Id = item.Id.VideoId
	}

	if item.Id.Kind == KindPlaylist {
		// Request the items from the playlists
		videos, err := client.PlaylistItems.List(snippetParam).PlaylistId(result.Id).MaxResults(maxPlaylistResults).Do()

		if err != nil {
			println("Found no videos in playlist " + result.Title)

			return nil, errors.New("found no videos in playlist")
		}

		results := make([]SearchResult, 0)

		// Parse results in the playlist
		for _, item := range videos.Items {
			s := item.Snippet

			vr := SearchResult{
				Id:    s.ResourceId.VideoId,
				Title: s.Title,
				Kind:  KindVideo,
			}

			results = append(results, vr)
		}

		result.Videos = results
	}

	return &result, nil
}

// VideosFromPlaylist returns all the videoIds associated with a playlist
func VideosFromPlaylist(id string) (*SearchResult, error) {
	// Find the playlist details
	playlists, err := client.Playlists.List(snippetParam).Id(id).Do()

	// Nothing found for the playlist id
	if err != nil || len(playlists.Items) < 1 {
		println("found nothing for playlists " + id)

		return nil, errors.New("failed to retrieve playlists")
	}

	// Set the playlist title
	title := playlists.Items[firstResult].Snippet.Title

	// Request the items from the playlists
	videos, err := client.PlaylistItems.List(snippetParam).PlaylistId(id).MaxResults(maxPlaylistResults).Do()

	if err != nil || len(videos.Items) == 0 {
		println("Found no videos in playlist " + title)

		return nil, errors.New("found no videos in playlist")
	}

	// Slice to hold the videos of the playlist
	results := make([]SearchResult, 0)

	// Parse results in the playlist
	for _, item := range videos.Items {
		s := item.Snippet

		vr := SearchResult{
			Id:    s.ResourceId.VideoId,
			Title: s.Title,
			Kind:  KindVideo,
		}

		results = append(results, vr)
	}

	// Include the playlist title with the videos
	return &SearchResult{
		Title:  title,
		Id:     id,
		Kind:   KindPlaylist,
		Videos: results,
	}, nil
}
