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
  ConfigPlayPathEnding = "RobocraftX_Data/StreamingAssets/Games/Freejam"
  ConfigIconPath = "icon.svg"
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
      config.BuildPath = filepath.FromSlash(os.Getenv("APPDATA")+"/../LocalLow/Freejam/RobocraftX/Games")
      config.PlayPath = filepath.FromSlash("C:/Program Files (x86)/Steam/steamapps/common/RobocraftX/"+ConfigPlayPathEnding)
    } else if runtime.GOOS == "linux" {
      config.BuildPath = filepath.FromSlash("~/.local/share/Steam/steamapps/compatdata/1078000/pfx/drive_c/users/steamuser/AppData/LocalLow/Freejam/RobocraftX/Games")
      config.PlayPath = filepath.FromSlash("~/.local/share/Steam/steamapps/common/RobocraftX/"+ConfigPlayPathEnding)
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
    config.IconPath = ConfigIconPath
  } else {
    data, _ := ioutil.ReadAll(file)
    json.Unmarshal(data, &config)
  }
  ForceGameCreatorTo = config.Creator
  ForceGameCreator = config.ForceCreator
  if config.DefaultSaveFolder != "" {
    DefaultSaveFolder = config.DefaultSaveFolder
  }
  if config.IconPath != "" {
    IconPath = config.IconPath
  }
  f, _ := os.Create(config.LogPath)
  log.Println("Log directed to "+config.LogPath)
  log.SetOutput(f)
  log.Println("Init complete")
}

func main() {
  var exitVal int
  log.Println("Starting main routine")
  config.Save()
  log.Println("RobocraftX Play Path: "+config.PlayPath)
  log.Println("RobocraftX Build Path: "+config.BuildPath)
  saveHandler := NewSaveHandler(config.PlayPath, config.BuildPath)
  activeDisplay = NewDisplay(saveHandler)
  activeDisplay.Start()
  exitVal, _ = activeDisplay.Join()
  if exitVal == 20 { // set new install dir
    log.Println("Display requested an update to PlayPath")
    config.PlayPath = filepath.FromSlash(NewInstallPath+"/"+ConfigPlayPathEnding)
    log.Println("New RobocraftX Play Path: "+config.PlayPath)
    config.Save()
    exitVal = 0
  }
  log.Println("rxsm terminated")
  os.Exit(exitVal) // this prevents defered operations, which may cause issues
}

// start of Config
type Config struct {
  PlayPath string `json:"play-path"`
  BuildPath string `json:"build-path"`
  Creator string `json:"creator"`
  ForceCreator bool `json:"force-creator?"`
  LogPath string `json:"log"`
  DefaultSaveFolder string `json:"copyable-save"`
  IconPath string `json:"icon"`
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
