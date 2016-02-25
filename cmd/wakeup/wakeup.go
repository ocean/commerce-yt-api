package main

import (
  "os"
  "log"
  "net/http"
)

func main() {
  url := os.Getenv("WAKE_UP_URL")
  if url == "" {
    log.Fatal("$WAKE_UP_URL must be set.")
  }
  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()
  log.Printf("Mmmf, what? Ok, waking up. Status: %s Protocol: %s", resp.Status, resp.Proto)
}
