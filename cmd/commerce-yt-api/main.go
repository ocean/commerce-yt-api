/*
Package main implements an API for proxying YouTube API requests.
Endpoints mirror their YouTube counterparts as much as possible,
so as little change as possible is required in the client (Drupal) codebase.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// searchRequest
//
// Simple search request, used for the SCALD_YOUTUBE_API Search request endpoint.
// https://www.googleapis.com/youtube/v3
// + /search?key=' . $api_key . '&q=' . $q . '&part=snippet&order=rating&type=video,playlist
func searchRequest(c *gin.Context) {
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
	if err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	var out = string(body[:])
	c.String(http.StatusOK, out)
}

// rssFeedRequest
//
// RSS feed request, used for the SCALD_YOUTUBE_API RSS Feed request endpoint.
// https://www.googleapis.com/youtube/v3
// + /videos?id=' . $id . '&key=' . $api_key . '&part=snippet
func rssFeedRequest(c *gin.Context) {
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
	if err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	var out = string(body[:])
	c.String(http.StatusOK, out)
}

// watchRequest
//
// Used for the SCALD_YOUTUBE_WEB request endpoint.
// https://www.youtube.com/watch
// + /watch?v=' . $id
func watchRequest(c *gin.Context) {
	id := c.Query("v")
	log.Printf("video id = %s", id)
	resp, err := http.Get("https://www.youtube.com/watch?v=" + id)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	var out = string(body[:])
	c.String(http.StatusOK, out)
}

// thumbnailRequest
//
// Used for the SCALD_YOUTUBE_THUMBNAIL request endpoint.
// https://i.ytimg.com
func thumbnailRequest(c *gin.Context) {
	q := c.Query("q")
	log.Printf("query url = %s", q)
	resp, err := http.Get(q)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: add in content type checking
	if strings.HasSuffix(q, "jpg") {
		c.Data(http.StatusOK, "image/jpeg", body)
	} else if strings.HasSuffix(q, "png") {
		c.Data(http.StatusOK, "image/png", body)
	} else {
		c.String(http.StatusForbidden, "403 Forbidden: Image requests only.")
	}
}

// ping
//
// Return ping requests with a nice timestamp.
func ping(c *gin.Context) {
	var resp struct {
		Response  string    `json:"response"`
		Timestamp time.Time `json:"timestamp"`
	}
	resp.Response = "pong"
	resp.Timestamp = time.Now().Local()
	c.JSON(http.StatusOK, resp)
}

// home
//
// Simple front page, using a template for fun.
func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":   "Hi there!",
		"heading": "Welcome",
		"content": "... to the API.",
	})
}

// Set up required variables
var (
	port = 5000
)

// Start the main function
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

	// Call in the gin router
	router := gin.Default()

	// Serve the damn favicon.ico
	router.StaticFile("/favicon.ico", "./public/favicon.ico")

	// Simple front page, using a template for fun.
	router.LoadHTMLGlob("templates/*")
	router.GET("/", home)

	// Return ping requests with a nice timestamp.
	router.GET("/ping", ping)

	// ----- ACTUAL REAL THINGS

	// SCALD_YOUTUBE_API Search request
	// https://www.googleapis.com/youtube/v3
	// + /search?key=' . $api_key . '&q=' . $q . '&part=snippet&order=rating&type=video,playlist
	router.GET("/v1/search", searchRequest)

	// SCALD_YOUTUBE_API RSS Feed request
	// https://www.googleapis.com/youtube/v3
	// + /videos?id=' . $id . '&key=' . $api_key . '&part=snippet
	router.GET("/v1/videos", rssFeedRequest)

	// SCALD_YOUTUBE_WEB request
	// https://www.youtube.com/watch
	// + /watch?v=' . $id
	router.GET("/v1/watch", watchRequest)

	// SCALD_YOUTUBE_THUMBNAIL request
	// https://i.ytimg.com
	router.GET("/v1/thumbnail", thumbnailRequest)

	// ----- SOME TEST THINGS
	router.GET("/form-submissions", func(c *gin.Context) {
		formsAPIToken := os.Getenv("FORMS_API_TOKEN")
		if formsAPIToken == "" {
			log.Fatal("$FORMS_API_TOKEN must be set")
		}
		resp, err := http.Get("http://forms.commerce.wa.gov.au/api/forms/results?token=" + formsAPIToken)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		var out = string(body[:])
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
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%s", body)

		c.Header("Content-Type", "application/json; charset=utf-8")
		var out = string(body[:])
		c.String(http.StatusOK, out)
	})

	// Run, collaborate and listen.
	router.Run(":" + port)
}
