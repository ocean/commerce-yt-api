package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

var (
	port = 5000
)

func main() {

	var port string
	if os.Getenv("PORT") == "" {
		port = "5000"
	} else {
		port = os.Getenv("PORT")
	}

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()

  router.LoadHTMLGlob("templates/*")
  router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
      "title": "Hi there!",
      "heading": "Welcome",
      "content": "... to the API.",
    })
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// ----- ACTUAL REAL THINGS

	// SCALD_YOUTUBE_API Search request
	// https://www.googleapis.com/youtube/v3
	// + /search?key=' . $api_key . '&q=' . $q . '&part=snippet&order=rating&type=video,playlist
	router.GET("/v1/search", func(c *gin.Context) {
		key := c.Query("key")
		q := url.QueryEscape(c.Query("q"))
		suffix := "&part=snippet&order=rating&type=video,playlist"
		log.Printf("search query = %s", q)
		resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?key=%s&q=%s%s", key, q, suffix))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		// fmt.Printf("%s", body)

		var jsonContentType = []string{"application/json; charset=utf-8"}
		writeContentType(c.Writer, jsonContentType)
		var out string = string(body[:])
		c.String(http.StatusOK, out)
	})

	// SCALD_YOUTUBE_API RSS Feed request
	// https://www.googleapis.com/youtube/v3
	// + /videos?id=' . $id . '&key=' . $api_key . '&part=snippet
	router.GET("/v1/videos", func(c *gin.Context) {
		id := c.Query("id")
		key := c.Query("key")
		suffix := "&part=snippet"
		log.Printf("video id = %s", id)
		resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s%s", id, key, suffix))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		// fmt.Printf("%s", body)

		var jsonContentType = []string{"application/json; charset=utf-8"}
		writeContentType(c.Writer, jsonContentType)
		var out string = string(body[:])
		c.String(http.StatusOK, out)
	})

	// SCALD_YOUTUBE_WEB request
	// https://www.youtube.com/watch
	// + /watch?v=' . $id
	router.GET("/v1/watch", func(c *gin.Context) {
		id := c.Query("v")
		log.Printf("video id = %s", id)
		resp, err := http.Get("https://www.youtube.com/watch?v=" + id)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		var htmlContentType = []string{"text/html; charset=utf-8"}
		writeContentType(c.Writer, htmlContentType)
		var out string = string(body[:])
		c.String(http.StatusOK, out)
	})

	// SCALD_YOUTUBE_THUMBNAIL request
	// https://i.ytimg.com
	router.GET("/v1/thumbnail", func(c *gin.Context) {
		q := c.Query("q")
		log.Printf("query url = %s", q)
		resp, err := http.Get(q)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		var jpegContentType = []string{"image/jpeg"}
		writeContentType(c.Writer, jpegContentType)
		var out string = string(body[:])
		c.String(http.StatusOK, out)
	})

	// ----- MANY TEST THINGS
	router.GET("/form-submissions", func(c *gin.Context) {
		resp, err := http.Get("http://forms.commerce.wa.gov.au/api/forms/results?token=ZuesbwqGhQMTxTbytbj7qrBWR_E84lTCSYLiVL1yk8Q")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		var jsonContentType = []string{"application/json; charset=utf-8"}
		writeContentType(c.Writer, jsonContentType)
		var out string = string(body[:])
		c.String(http.StatusOK, out)
	})

	router.GET("/fuel/:suburb", func(c *gin.Context) {
		suburb := c.Param("suburb")
		resp, err := http.Get("http://nfwws.herokuapp.com/v1/s/" + suburb)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		// fmt.Printf("%s", body)

		var jsonContentType = []string{"application/json; charset=utf-8"}
		writeContentType(c.Writer, jsonContentType)
		var out string = string(body[:])
		c.String(http.StatusOK, out)
	})

	router.Run(":" + port)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
