//#!/bin/go

// Created 2019-07-26 by NGnius

package main

import (
  "os"
  "encoding/json"
  "fmt"
)

func init() {
  // load config file
}

func main() {
  fmt.Println("Hello world")
}

type Config struct {
  BasePath string `json:"installPath"`
  Creator string `json:"creator"`
}

func (c Config) save() (error) {
  file, openErr := os.Create("config.json")
  if openErr != nil {
    return openErr
  }
  out, marshalErr := json.Marshal(c)
  if marshalErr != nil {
    return marshalErr
  }
  file.Write(out)
  file.Sync()
  file.Close()
  return nil
}
