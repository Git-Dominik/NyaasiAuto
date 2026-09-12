package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	if err := qBitAuth(&client); err != nil {
		fmt.Println(err)
		return
	}

	startServer(&client)

}

func qBitAuth(client *http.Client) error {
	env, err := loadEnv()
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := client.PostForm("http://localhost:8080/api/v2/auth/login", url.Values{
		"username": {env.User},
		"password": {env.Pass},
	})
	if err != nil {
		fmt.Print(err)
		return err
	}

	fmt.Println(resp.Status)
	defer resp.Body.Close()
	return nil
}

func addTorrent(magnet string, client *http.Client) error {
	addResp, err := client.PostForm("http://localhost:8080/api/v2/torrents/add", url.Values{
		"urls": {magnet},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(addResp.Status)
	defer addResp.Body.Close()
	return nil
}

func startServer(client *http.Client) {
	router := gin.Default()

	router.POST("/add-torrent", func(c *gin.Context) {
		magnet := c.PostForm("magnet")
		if magnet == "" {
			c.String(http.StatusBadRequest, "Missing magnet parameneter")
			return
		}

		err := addTorrent(magnet, client)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "Failed to submit torrent")
			return
		}
		c.String(http.StatusOK, magnet)
	})

	router.Run(":8081")
}
