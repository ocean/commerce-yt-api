package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var (
	port = 5000
)

// SCALD_YOUTUBE_API Search request
// https://www.googleapis.com/youtube/v3
// + /search?key=' . $api_key . '&q=' . $q . '&part=snippet&order=rating&type=video,playlist
func searchRequest(c echo.Context) error {
	key := c.QueryParam("key")
	q := url.QueryEscape(c.QueryParam("q"))
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

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	var out = string(body[:])
	return c.String(http.StatusOK, out)
}

// SCALD_YOUTUBE_API RSS Feed request
// https://www.googleapis.com/youtube/v3
// + /videos?id=' . $id . '&key=' . $api_key . '&part=snippet
func rssFeedRequest(c echo.Context) error {
	id := c.QueryParam("id")
	key := c.QueryParam("key")
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

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	var out = string(body[:])
	return c.String(http.StatusOK, out)
}

// SCALD_YOUTUBE_WEB request
// https://www.youtube.com/watch
// + /watch?v=' . $id
func watchRequest(c echo.Context) error {
	id := c.QueryParam("v")
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

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	var out = string(body[:])
	return c.String(http.StatusOK, out)
}

// SCALD_YOUTUBE_THUMBNAIL request
// https://i.ytimg.com
// func thumbnailRequest(c echo.Context) error {
// 	q := c.QueryParam("q")
// 	log.Printf("query url = %s", q)
// 	resp, err := http.Get(q)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// TODO: add in content type checking
// 	if strings.HasSuffix(q, "jpg") {
// 		c.Response().Header().Set(echo.HeaderContentType, echo.ContentTypeByExtension(q))
// 		// c.Data(http.StatusOK, "image/jpeg", body)
// 		c.Response().WriteHeader(http.StatusOK)
// 		_, err = c.Response().Write(body)
// 		// return c.JSONBlob(http.StatusOK, body)
// 	} else if strings.HasSuffix(q, "png") {
// 		c.Response().Header().Set(echo.HeaderContentType, echo.ContentTypeByExtension(q))
// 		// c.Data(http.StatusOK, "image/png", body)
// 		return c.JSONBlob(http.StatusOK, body)
// 	} else {
// 		return c.String(http.StatusForbidden, "403 Forbidden: Image requests only.")
// 	}
// 	return c.Render()
// }

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

	e := echo.New()

	e.File("/", "public/index.html")

	e.GET("/ping", func(c echo.Context) error {
		var resp struct {
			Response  string    `json:"response"`
			Timestamp time.Time `json:"timestamp"`
		}
		resp.Response = "pong"
		resp.Timestamp = time.Now().Local()
		return c.JSON(http.StatusOK, resp)
	})

	// ----- ACTUAL REAL THINGS

	// SCALD_YOUTUBE_API Search request
	// https://www.googleapis.com/youtube/v3
	// + /search?key=' . $api_key . '&q=' . $q . '&part=snippet&order=rating&type=video,playlist
	e.GET("/v1/search", searchRequest)

	// SCALD_YOUTUBE_API RSS Feed request
	// https://www.googleapis.com/youtube/v3
	// + /videos?id=' . $id . '&key=' . $api_key . '&part=snippet
	e.GET("/v1/videos", rssFeedRequest)

	// SCALD_YOUTUBE_WEB request
	// https://www.youtube.com/watch
	// + /watch?v=' . $id
	e.GET("/v1/watch", watchRequest)

	// SCALD_YOUTUBE_THUMBNAIL request
	// https://i.ytimg.com
	// e.GET("/v1/thumbnail", thumbnailRequest)

	e.Run(standard.New(":" + port))
}
