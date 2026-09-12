package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var episodeRegex = regexp.MustCompile(`(?i)(?:S(\d{1,2}))?\s*(?:E|Episode|\s+-\s+)\s*(\d{2,3})`)

func extractEpisode(title string) (int, int, bool) {
	matches := episodeRegex.FindStringSubmatch(title)

	if len(matches) < 3 {
		return 0, 0, false
	}

	season := 1
	if matches[1] != "" {
		season, _ = strconv.Atoi(matches[1])
	}

	episode, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, false
	}

	return season, episode, true
}

func buildNyaaURL(animeTitle string, uploader string) string {
	cleanTitle := strings.ReplaceAll(animeTitle, ":", "")

	searchQuery := fmt.Sprintf("[%s] %s", uploader, cleanTitle)

	baseURL := "https://nyaa.si/?page=rss&c=0_0&f=0&q="
	return baseURL + url.QueryEscape(searchQuery)
}
