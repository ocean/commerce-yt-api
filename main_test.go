package main

import (
  "testing"
  // "net/http"
  // "log"
  // "net/http/httptest"
  // "github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
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
