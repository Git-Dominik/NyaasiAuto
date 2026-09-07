package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/mmcdole/gofeed"
)

type Progress struct {
	Season  int
	Episode int
}

type Target struct {
	Title      string
	Uploaders  []string
	latestSeen Progress
}

var episodeRegex = regexp.MustCompile(`(?i)(?:S(\d{1,2}))?\s*(?:E|Episode|\s+-\s+)\s*(\d{2,3})`)

func main() {
	fp := gofeed.NewParser()

	targets := []Target{
		{
			Title:      "rezero",
			Uploaders:  []string{"ToonsHub", "FBI"},
			latestSeen: Progress{Season: 4, Episode: 10},
		},
		{
			Title:      "mushoku tensei",
			Uploaders:  []string{"Feibanyama"},
			latestSeen: Progress{Season: 3, Episode: 10},
		},
	}

	for _, target := range targets {

		rssURL := buildNyaaURL(target.Title, target.Uploaders)

		feed, err := fp.ParseURL(rssURL)
		if err != nil {
			fmt.Println("Error fetching feed for", target.Title, ":", err)
			continue
		}

		for _, item := range feed.Items {

			season, episode, found := extractEpisode(item.Title)
			if !found {
				fmt.Println("[UNMATCHED]", item.Title)
				continue
			}

			episodeTag := fmt.Sprintf("S%02dE%02d", season, episode)
			isNewer := season > target.latestSeen.Season || (season == target.latestSeen.Season && episode > target.latestSeen.Episode)

			if isNewer {
				fmt.Println("[NEW EPISODE]", episodeTag, item.Title)
				fmt.Println(item.Link)
			} else {
				fmt.Println("[SKIP]", episodeTag, item.Title)
			}
		}
	}

}

func buildNyaaURL(animeTitle string, uploaders []string) string {
	var queryParts []string

	cleanTitle := strings.ReplaceAll(animeTitle, ":", "")

	if len(uploaders) > 0 {
		queryParts = append(queryParts, "("+strings.Join(uploaders, "|")+")")
	}

	queryParts = append(queryParts, cleanTitle)
	searchQuery := strings.Join(queryParts, " ")

	baseURL := "https://nyaa.si/?page=rss&c=0_0&f=0&q="
	return baseURL + url.QueryEscape(searchQuery)
}

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
