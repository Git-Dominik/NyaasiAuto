package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://nyaa.si/?page=rss&q=%5BFeibanyama%5D+Mushoku+Tensei+Jobless+Reincarnation&c=0_0&f=0")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(feed.Title)

	for _, i := range feed.Items {
		fmt.Println(i.Title, "|", i.Link, "|", i.Published, "|", i.GUID)
	}
}
