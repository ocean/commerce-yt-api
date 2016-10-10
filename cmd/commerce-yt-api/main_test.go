package main

import (
	"testing"

	"github.com/gin-gonic/gin"

	"net/http"
	"net/http/httptest"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHomePage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.LoadHTMLGlob("../../templates/*")
	router.GET("/", home)

	Convey("When you hit the root URL", t, func() {
		Convey("The result should contain Welcome", func() {
			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			// log.Println(resp)
			So(resp.Code, ShouldHaveSameTypeAs, 1)
			So(resp.Code, ShouldEqual, 200)
			So(resp.HeaderMap, ShouldContainKey, "Content-Type")
			So(resp.HeaderMap.Get("Content-Type"), ShouldContainSubstring, "text/html")
			So(resp.Body.String(), ShouldContainSubstring, "Welcome")
		})
	})
}

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/ping", ping)

	Convey("When you hit the ping URL", t, func() {
		Convey("The result should contain pong", func() {
			req, _ := http.NewRequest("GET", "/ping", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			// log.Println(resp)
			So(resp.Code, ShouldHaveSameTypeAs, 1)
			So(resp.Code, ShouldEqual, 200)
			So(resp.HeaderMap, ShouldContainKey, "Content-Type")
			So(resp.HeaderMap.Get("Content-Type"), ShouldContainSubstring, "application/json")
			So(resp.Body.String(), ShouldContainSubstring, "pong")
		})
	})
}

func TestWatchRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/v1/watch", watchRequest)

	Convey("When you hit the watch URL with the video id for PSY - GANGNAM STYLE", t, func() {
		Convey("The result should contain PSY - GANGNAM", func() {
			req, _ := http.NewRequest("GET", "/v1/watch?v=9bZkp7q19f0", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			// log.Println(resp)
			// log.Println(resp.Status())
			// log.Println(resp.HeaderMap["Content-Type"])
			So(resp.Code, ShouldHaveSameTypeAs, 1)
			So(resp.Code, ShouldEqual, 200)
			So(resp.HeaderMap, ShouldContainKey, "Content-Type")
			So(resp.HeaderMap.Get("Content-Type"), ShouldContainSubstring, "text/html")
			So(resp.Body.String(), ShouldContainSubstring, "PSY - GANGNAM")
		})
	})
}

func TestThumbnailRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/v1/thumbnail", thumbnailRequest)

	Convey("When you hit the thumbnail URL for PSY - GANGNAM STYLE", t, func() {
		Convey("The result should be an image", func() {
			req, _ := http.NewRequest("GET", "/v1/thumbnail?q=https://i.ytimg.com/vi/9bZkp7q19f0/hqdefault.jpg", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			So(resp.Code, ShouldHaveSameTypeAs, 1)
			So(resp.Code, ShouldEqual, 200)
			So(resp.HeaderMap, ShouldContainKey, "Content-Type")
			So(resp.HeaderMap.Get("Content-Type"), ShouldContainSubstring, "image")
		})
	})
}

func TestBadThumbnailRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/v1/thumbnail", thumbnailRequest)

	Convey("When you hit the thumbnail URL looking for some random thing", t, func() {
		Convey("The result should be a '403 Forbidden' error", func() {
			req, _ := http.NewRequest("GET", "/v1/thumbnail?q=https://www.google.com.au", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			So(resp.Code, ShouldHaveSameTypeAs, 1)
			So(resp.Code, ShouldEqual, 403)
			So(resp.HeaderMap, ShouldContainKey, "Content-Type")
			So(resp.HeaderMap.Get("Content-Type"), ShouldContainSubstring, "text/plain")
		})
	})
}
