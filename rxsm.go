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
)

const (
  // config defaults
  ConfigCreator = "unknown"
  ConfigLogPath = "rxsm.log"
  ConfigSaveFolder = "default_save"
)

var configPath string = "config.json"
var config Config

var activeDisplay IDisplayGoroutine

func init() {
  log.Println("Starting init")
  // load config file
  file, openErr := os.Open(configPath)
  if openErr != nil {
    // no config file found, use default config
    if runtime.GOOS == "windows" {
      config.PlayPath = filepath.FromSlash("%APPDATA%/../LocalLow/Freejam/RobocraftX/Games")
      config.BuildPath = filepath.FromSlash("C:/Program Files (x86)/Steam/steamapps/common/RobocraftX/RobocraftX_Data/StreamingAssets/Games/Freejam")
    } else if runtime.GOOS == "linux" {
      config.PlayPath = filepath.FromSlash("~/.local/share/Steam/steamapps/compatdata/1078000/pfx/drive_c/users/steamuser/AppData/LocalLow/Freejam/RobocraftX/Games")
      config.BuildPath = filepath.FromSlash("~/.local/share/Steam/steamapps/common/RobocraftX/RobocraftX_Data/StreamingAssets/Games/Freejam")
    } else if runtime.GOOS == "darwin" { // macOS
      // support doesn't really matter until SteamPlay or FJ supports MacOS
      log.Fatal("OS detected as macOS (unsupported)")
    } else {
      log.Fatal("No default config for OS: "+runtime.GOOS)
    }
    config.Creator = ConfigCreator
    config.LogPath = ConfigLogPath
    config.ForceCreator = false
    config.DefaultSaveFolder = ConfigSaveFolder
  } else {
    data, _ := ioutil.ReadAll(file)
    json.Unmarshal(data, &config)
  }
  ForceGameCreatorTo = config.Creator
  ForceGameCreator = config.ForceCreator
  if config.DefaultSaveFolder != "" {
    DefaultSaveFolder = config.DefaultSaveFolder
  }
  f, _ := os.Create(config.LogPath)
  log.Println("Log directed to "+config.LogPath)
  log.SetOutput(f)
  log.Println("Init complete")
}

func main() {
  log.Println("Starting main routine")
  config.Save()
  log.Println("RobocraftX Play Path: "+config.PlayPath)
  log.Println("RobocraftX Build Path: "+config.BuildPath)
  saveHandler := NewSaveHandler(config.PlayPath, config.BuildPath)
  activeDisplay = NewDisplay(saveHandler)
  activeDisplay.Start()
  activeDisplay.Join()
  log.Println("rxsm terminated")
}

// start of Config
type Config struct {
  PlayPath string `json:"play-path"`
  BuildPath string `json:"build-path"`
  Creator string `json:"creator"`
  ForceCreator bool `json:"force-creator?"`
  LogPath string `json:"log"`
  DefaultSaveFolder string `json:"copyable-save"`
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
