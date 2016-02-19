package main

import (
  "log"
  "net/http"
  "os"
  // "fmt"
  // "encoding/json"
  "io/ioutil"
  "github.com/gin-gonic/gin"
)

var (
  port = 5000
)

func main() {

  var port string
  if (os.Getenv("PORT") == "") {
    port = "5000"
  } else {
    port = os.Getenv("PORT")
  }

  if port == "" {
		log.Fatal("$PORT must be set")
	}

  router := gin.Default()

  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
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
    var out string =  string(body[:])
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
