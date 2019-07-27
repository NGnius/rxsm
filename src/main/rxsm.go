//#!/bin/go

// Created 2019-07-26 by NGnius

package main

import (
  "os"
  "encoding/json"
  "path/filepath"
  "runtime"
  "io/ioutil"
  "log"
  //"fmt"

  "./saver"
)

const (
  // config defaults
  ConfigCreator = "unknown"
  ConfigLogPath = "log.txt"
)

var configPath string = "config.json"
var config Config

func init() {
  log.Println("Starting init")
  // load config file
  file, openErr := os.Open(configPath)
  if openErr != nil {
    // no config file found, use default config
    if runtime.GOOS == "windows" {
      config.BasePath = filepath.FromSlash("C:/Program Files (x86)/Steam/steamapps/common/RobocraftX")
    } else if runtime.GOOS == "linux" {
      config.BasePath = filepath.FromSlash("/home/ngnius/.local/share/Steam/steamapps/common/RobocraftX")
    } else if runtime.GOOS == "darwin" { // macOS
      // support doesn't really matter until SteamPlay or FJ supports MacOS
    } else {
      log.Println("No default config for OS: "+runtime.GOOS)
    }
    config.Creator = ConfigCreator
    config.LogPath = ConfigLogPath
    config.ForceCreator = false
  } else {
    data, _ := ioutil.ReadAll(file)
    json.Unmarshal(data, &config)
  }
  saver.ForceGameCreatorTo = config.Creator
  saver.ForceGameCreator = config.ForceCreator
  f, _ := os.Create(config.LogPath)
  log.Println("Log directed to "+config.LogPath)
  log.SetOutput(f)
  log.Println("Init complete")
}

func main() {
  log.Println("Starting main routine")
  config.Save()
  log.Println("RobocraftX Install Path: "+config.BasePath)
  saveHandler := saver.NewSaveHandler(config.BasePath)
  log.Println("rxsm terminated")
}

// start of Config
type Config struct {
  BasePath string `json:"installPath"`
  Creator string `json:"creator"`
  ForceCreator bool `json:"force-creator?"`
  LogPath string `json:"log"`
}

func (c Config) Save() (error) {
  file, openErr := os.Create(configPath)
  if openErr != nil {
    return openErr
  }
  out, marshalErr := json.MarshalIndent(c, "", "  ")
  if marshalErr != nil {
    return marshalErr
  }
  file.Write(out)
  file.Sync()
  file.Close()
  return nil
}
// end of Config
