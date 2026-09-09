package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/superturkey650/go-qbittorrent/qbt"

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

func checkUpdate(configFilePath string) {
	targets, err := loadTargets(configFilePath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	fp := gofeed.NewParser()

	for i := range targets {
		foundNewEpisode := false

		for _, uploader := range targets[i].Uploaders {
			if foundNewEpisode {
				break
			}
			rssURL := buildNyaaURL(targets[i].Title, uploader)
			feed, err := fp.ParseURL(rssURL)
			if err != nil {
				continue
			}

			for {
				expectedSeason := targets[i].LatestSeen.Season
				expectedEpisode := targets[i].LatestSeen.Episode + 1
				foundThisPass := false

				for _, item := range feed.Items {
					season, episode, found := extractEpisode(item.Title)
					if !found {
						continue
					}

					if season == expectedSeason && episode == expectedEpisode {
						if err := torrentLinks(item.Link); err != nil {
							fmt.Print(err)
						} else {
							episodeTag := fmt.Sprintf("S%02dE%02d", season, episode)
							fmt.Println("[NEW EPISODE]", episodeTag, item.Title)
							fmt.Println(item.Link)

							targets[i].LatestSeen.Season = season
							targets[i].LatestSeen.Episode = episode
							targets[i].saveProgress(configFilePath, targets)
							fmt.Println(targets[i].LatestSeen)
							foundNewEpisode = true
							foundThisPass = true
							break
						}

					}

				}

				if !foundThisPass {
					break
				}
			}

		}

	}
}
func main() {
	configFilePath := "config.json"

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			checkUpdate(configFilePath)
			fmt.Println("check")
		}
	}()

	wg.Wait()
}

func torrentLinks(torrent string) error {
	qb := qbt.NewClient("http://localhost:8080/")
	if err := qb.Login("admin", "your-secret-password"); err != nil {
		return err
	}

	options := qbt.DownloadOptions{}

	if err := qb.DownloadLinks([]string{torrent}, options); err != nil {
		return err
	}

	return nil
}

func (t *Target) saveProgress(configFilePath string, targets []Target) error {
	updateJson, err := json.MarshalIndent(targets, "", " ")
	if err != nil {
		panic(err)
	}

	return os.WriteFile(configFilePath, updateJson, 0644)

}

func loadTargets(configFilePath string) ([]Target, error) {
	file, err := os.ReadFile(configFilePath)
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
