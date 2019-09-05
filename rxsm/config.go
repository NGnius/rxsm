// Created 2019-09-04 by NGnius

package main

import (
  "os"
  "runtime"
  "path/filepath"
  "encoding/json"
  "io/ioutil"
  "log"
)

const (
  ConfigPlayPathEnding = "RobocraftX_Data/StreamingAssets/Games/Freejam"
)

var globalConfigPath string = "config.json"
var GlobalConfig *Config

// start of Config

type Config struct {
  PlayPath string `json:"play-path"`
  BuildPath string `json:"build-path"`
  Creator string `json:"creator"`
  ForceCreator bool `json:"force-creator?"`
  LogPath string `json:"log"`
  DefaultSaveFolder string `json:"copyable-save"`
  IconPath string `json:"icon"`
  ForceUniqueIds bool `json:"force-unique-ids?"`
  SettingsIconPath string `json:"settings-icon"`
  Version string `json:"version"`
  lastVersion string
  path string
}

func DefaultConfig() (c *Config) {
  c = &Config{}
  c.Creator = "unknown"
  c.ForceCreator = false
  c.LogPath = "rxsm.log"
  c.DefaultSaveFolder = "default_save"
  c.IconPath = "icon.svg"
  c.ForceUniqueIds = false
  c.SettingsIconPath = "settings.svg"
  if runtime.GOOS == "windows" {
    c.BuildPath = filepath.FromSlash(os.Getenv("APPDATA")+"/../LocalLow/Freejam/RobocraftX/Games")
    c.PlayPath = filepath.FromSlash("C:/Program Files (x86)/Steam/steamapps/common/RobocraftX/"+ConfigPlayPathEnding)
  } else if runtime.GOOS == "linux" {
    c.BuildPath = filepath.FromSlash("~/.local/share/Steam/steamapps/compatdata/1078000/pfx/drive_c/users/steamuser/AppData/LocalLow/Freejam/RobocraftX/Games")
    c.PlayPath = filepath.FromSlash("~/.local/share/Steam/steamapps/common/RobocraftX/"+ConfigPlayPathEnding)
  } else if runtime.GOOS == "darwin" { // macOS
    // support doesn't really matter until SteamPlay or FJ supports MacOS
    log.Fatal("OS detected as macOS (unsupported)")
  } else {
    log.Fatal("No default config for OS: "+runtime.GOOS)
  }
  return
}

func (c *Config) Save() (error) {
  file, openErr := os.Create(c.path)
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

func (c *Config) LastVersion() (string) {
  return c.lastVersion
}

func (c *Config) load(path string) (error) {
  file, openErr := os.Open(path)
  if openErr != nil {
    return openErr
  }
  data, readErr := ioutil.ReadAll(file)
  if readErr != nil {
    return readErr
  }
  unmarshalErr := json.Unmarshal(data, c)
  if unmarshalErr != nil {
    return unmarshalErr
  }
  c.path = path
  // input cleaning
  c.PlayPath = filepath.FromSlash(c.PlayPath)
  c.BuildPath = filepath.FromSlash(c.BuildPath)
  c.LogPath = filepath.FromSlash(c.LogPath)
  c.IconPath = filepath.FromSlash(c.IconPath)
  c.SettingsIconPath = filepath.FromSlash(c.SettingsIconPath)
  c.lastVersion = c.Version
  c.Version = RXSMVersion
  return nil
}

// end of Config

func LoadGlobalConfig() (*Config, error) {
  c := DefaultConfig()
  err := c.load(globalConfigPath)
  if err == nil {
    GlobalConfig = c
  }
  return c, err
}
