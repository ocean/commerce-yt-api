package main

import (
	"testing"

	// "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"net/http/httptest"
)

func TestPing(t *testing.T) {
	// w := httptest.NewRecorder()
	//
	// router := gin.Default()
	//
	// req, err := http.NewRequest("GET", "http://localhost:5000/ping", nil)
	// if err != nil {
	//   log.Fatal(err)
	// }

	// reqerr := router.GET(req)

	// assert.NoError(t, reqerr)
	// assert.Equal(t, w.Body.String(), "{\"message\": \"pong\"}")
	assert.Equal(t, "ping", "ping")
}

func TestHomePage(t *testing.T)  {
	// h := httptest.NewRecorder()

	// req, err := http.NewRequest("GET", "http://localhost:5000/", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := getRouter(true)

	// r := gin.Default()
	r.GET("/ping", ping)

	//
	// req, err := http.NewRequest("GET", "http://localhost:5000/ping", nil)
	// if err != nil {
	//   log.Fatal(err)
	// }

	// reqerr := router.GET(req)

	// assert.NoError(t, reqerr)
	// assert.Equal(t, w.Body.String(), "{\"message\": \"pong\"}")

	Convey("When you hit the root URL", t, func() {
		Convey("The result should contain Welcome", func() {
			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			log.Println(resp)
			So(resp.Body.String(), ShouldContainSubstring, "Welcome")
		})
	})

	Convey("When you hit the ping URL", t, func() {
		Convey("The result should contain pong", func() {
			req, _ := http.NewRequest("GET", "/ping", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			log.Println(resp)
			So(resp.Body.String(), ShouldContainSubstring, "pong")
		})
	})
}
