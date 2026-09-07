package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mmcdole/gofeed"
)

type Progress struct {
	Season  int `json:"season"`
	Episode int `json:"episode"`
}

type Target struct {
	Title      string   `json:"title"`
	Uploaders  []string `json:"uploaders"`
	LatestSeen Progress `json:"latest_seen"`
}

var episodeRegex = regexp.MustCompile(`(?i)(?:S(\d{1,2}))?\s*(?:E|Episode|\s+-\s+)\s*(\d{2,3})`)

func main() {
	targets, err := loadTargets("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	fp := gofeed.NewParser()

	for _, target := range targets {
		foundNewEpisode := false

		for _, uploader := range target.Uploaders {
			if foundNewEpisode {
				break
			}
			rssURL := buildNyaaURL(target.Title, uploader)
			feed, err := fp.ParseURL(rssURL)
			if err != nil {
				continue
			}

			for _, item := range feed.Items {
				season, episode, found := extractEpisode(item.Title)
				if !found {
					continue
				}

				episodeTag := fmt.Sprintf("S%02dE%02d", season, episode)
				isNewer := season > target.LatestSeen.Season || (season == target.LatestSeen.Season && episode > target.LatestSeen.Episode)

				if isNewer {
					fmt.Println("[NEW EPISODE]", episodeTag, item.Title)
					fmt.Println(item.Link)

					foundNewEpisode = true
				} // else {
				// 	fmt.Println("[SKIP]", episodeTag, item.Title)
				// }
			}
		}

	}
}

func loadTargets(filename string) ([]Target, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var targets []Target

	if err := json.Unmarshal(file, &targets); err != nil {
		panic(err)
	}

	return targets, nil
}

func buildNyaaURL(animeTitle string, uploader string) string {
	cleanTitle := strings.ReplaceAll(animeTitle, ":", "")

	searchQuery := fmt.Sprintf("[%s] %s", uploader, cleanTitle)

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
