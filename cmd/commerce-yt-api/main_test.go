package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	// "github.com/stretchr/testify/assert"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"net/http/httptest"
)

func TestHomePage(t *testing.T)  {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/", home)

	Convey("When you hit the root URL", t, func() {
		Convey("The result should contain Welcome", func() {
			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			log.Println(resp)
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
			log.Println(resp)
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
			So(resp.Body.String(), ShouldContainSubstring, "PSY - GANGNAM")
		})
	})
}
