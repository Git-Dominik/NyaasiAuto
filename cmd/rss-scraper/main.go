package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

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
							log.Fatal(err)
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
	targets, _ := loadTargets(configFilePath)

	titles := make(map[string]int)
	for _, t := range targets {
		titles[t.Title] = t.LatestSeen.Episode + 1
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		defer wg.Done()
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			checkUpdate(configFilePath)
			fmt.Println("checking for", titles)
		}
	})

	wg.Wait()
}
