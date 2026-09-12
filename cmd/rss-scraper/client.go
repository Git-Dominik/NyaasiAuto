package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func torrentLinks(torrent string) error {
	resp, err := http.PostForm("http://localhost:8081/add-torrent", url.Values{
		"magnet": {torrent},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
