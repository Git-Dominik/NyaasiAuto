package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mmcdole/gofeed"
)

type Progress struct {
	Season  int
	Episode int
}

var episodeRegex = regexp.MustCompile(`(?i)(?:S(\d{1,2}))?\s*(?:E|Episode|\s+-\s+)\s*(\d{2,3})`)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://nyaa.si/?page=rss&q=%5BFeibanyama%5D+Mushoku+Tensei+Jobless+Reincarnation&c=0_0&f=0")
	if err != nil {
		fmt.Println(err)
	}

	latestSeen := Progress{
		Season:  4,
		Episode: 2,
	}

	for _, item := range feed.Items {
		season, episode, found := extractEpisode(item.Title)
		if !found {
			fmt.Println("[UNMATCHED]", item.Title)
			continue
		}

		episodeTag := fmt.Sprintf("S%02dE%02d", season, episode)
		isNewer := season > latestSeen.Season || (season == latestSeen.Season && episode > latestSeen.Episode)

		if isNewer {
			fmt.Println("[NEW EPISODE]", episodeTag, item.Title)
			fmt.Println(item.Link)
		} else {
			fmt.Println("[SKIP]", episodeTag, item.Title)
		}
	}
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
