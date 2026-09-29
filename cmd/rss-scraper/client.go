package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type torrentStatus struct {
	Name     string  `json:"name"`
	Progress float64 `json:"progress"`
}

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

func torrentInfo() ([]torrentStatus, error) {
	resp, err := http.Get("http://localhost:8081/info")
	if err != nil {
		return nil, fmt.Errorf("Failed to send get request to /info %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var all []torrentStatus
	if err := json.Unmarshal(body, &all); err != nil {
		return nil, fmt.Errorf("unmarshal succeeded: %w", err)
	}

	return all, nil
}

func torrentInfoMap() (map[string]float64, error) {
	list, err := torrentInfo()
	if err != nil {
		return nil, err
	}

	progressMap := make(map[string]float64)
	for _, status := range list {
		progressMap[status.Name] = status.Progress
	}
	return progressMap, nil
}
