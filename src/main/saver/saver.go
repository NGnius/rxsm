// Created 2019-07-26 by NGnius

package saver

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "encoding/json"
  "strconv"
  //"log"
)

const (
  Games = "RobocraftX_Data/StreamingAssets/Games"
  Play = "FreeJam"
  Build = "Player"
  GameStart = "Game_"

  GameDataFile = "GameData.json"
  GameSaveFile = "GameSave.RCX"
  ThumbnailFile = "Thumbnail.jpg"
)

var (
  ForceGameCreator = false
  ForceGameCreatorTo = ""
)

// start of SaveHandler
type SaveHandler struct {
  installPath string
  gamesPath string
  playPath string
  buildPath string
  PlaySaves []Save
  BuildSaves []Save
}

func NewSaveHandler(installPath string) (SaveHandler) {
  newSaveHandler := SaveHandler {
    installPath: installPath,
    gamesPath: filepath.Join(installPath, Games),
    playPath: filepath.Join(installPath, Games, Play),
    buildPath: filepath.Join(installPath, Games, Build)}
  newSaveHandler.PlaySaves = newSaveHandler.getSaves(newSaveHandler.playPath)
  newSaveHandler.BuildSaves = newSaveHandler.getSaves(newSaveHandler.buildPath)
  return newSaveHandler
}

func (sv SaveHandler) getSaves(saveFolder string) ([]Save){
  var saves []Save
  folders := getFoldersInFolder(saveFolder)
  for _, folder := range folders {
    s, sErr := NewSave(folder)
    if sErr == nil {
      saves = append(saves, s)
    }
  }
  return saves
}
// end of SaveHandler

// start Save
type Save struct {
  Data GameData
  dataPath string
  savePath string
  thumbnailPath string
  folder string
}

func NewSave(folder string) (Save, error) {
  newSave := Save {
    dataPath: filepath.Join(folder, GameDataFile),
    savePath: filepath.Join(folder, GameSaveFile),
    thumbnailPath: filepath.Join(folder, ThumbnailFile),
    folder: folder}
  newGD, gdErr := NewGameData(newSave.dataPath)
  newSave.Data = newGD
  if gdErr != nil {
    return newSave, gdErr
  }
  return newSave, nil
}

func (s *Save) Move(to string) (error) {
  moveErr := os.Rename(s.folder, to)
  if moveErr != nil {
    return moveErr
  }
  s.folder = to
  s.dataPath = filepath.Join(to, GameDataFile)
  s.Data.path = s.dataPath
  s.savePath = filepath.Join(to, GameSaveFile)
  s.thumbnailPath = filepath.Join(to, ThumbnailFile)
  return nil
}

func (s *Save) MoveToId() (error) {
  idDir, _ := filepath.Split(s.folder)
  idDir = filepath.Join(idDir, GameStart+doubleDigitStr(s.Data.Id))
  return s.Move(idDir)
}

func (s *Save) MoveToFirst() (error) {
  firstDir, _ := filepath.Split(s.folder)
  firstDir = filepath.Join(firstDir, GameStart+"00")
  return s.Move(firstDir)
}

func (s *Save) MoveOut() (error) {
  exists := true
  rootFolder, _ := filepath.Split(s.folder)
  lastDir := filepath.Join(rootFolder, GameStart+doubleDigitStr(s.Data.Id))
  if s.folder == lastDir && s.Data.Id != 0 { // don't move out if in zero-th place
    return nil
  }
  dirCheckLoop: for {
    if !exists {
      break dirCheckLoop
    }
    _, err := os.Stat(lastDir)
    exists = err == nil || os.IsExist(err)
    if exists {
      s.Data.Id += 1
      lastDir = filepath.Join(rootFolder, GameStart+doubleDigitStr(s.Data.Id))
      if s.Data.Id > 9999 { // sanity check
        return err
      }
    }
  }
  return s.Move(lastDir)
}
// end Save

// start of GameData
type GameData struct { // reference: GameData.json
  Id int `json:"GameID"`
  Name string `json:"GameName"`
  Description string `json:"GameDescription"`
  Creator string `json:"CreatorName"`
  path string
  isInited bool
}

func NewGameData(path string) (GameData, error) {
  var gd GameData
  f, openErr := os.Open(path)
  if openErr != nil {
    return gd, openErr
  }
  data, _ := ioutil.ReadAll(f)
  marshErr := json.Unmarshal(data, &gd)
  if marshErr != nil {
    return gd, marshErr
  }
  // check forced values
  if ForceGameCreator {
    gd.Creator = ForceGameCreatorTo
  }
  gd.path = path
  gd.isInited = true
  return gd, nil
}

func (gd GameData) Save() (error) {
  file, openErr := os.Create(gd.path)
  if openErr != nil {
    return openErr
  }
  out, marshalErr := json.MarshalIndent(gd, "", "  ")
  if marshalErr != nil {
    return marshalErr
  }
  file.Write(out)
  file.Sync()
  file.Close()
  return nil
}
// end of GameData

// helper functions

func getFoldersInFolder(dirpath string) []string {
  var folders []string
  dir, dirErr := os.Open(dirpath)
  if dirErr != nil {
    return folders
  }
  paths, readErr := dir.Readdirnames(0)
  if readErr != nil {
    return folders
  }
  for _, path := range paths {
    fullpath := filepath.Join(dirpath, path)
    info, statErr := os.Stat(fullpath)
    if statErr == nil && info.IsDir(){ // ignore errored pathes
      folders = append(folders, fullpath)
    }
  }
  return folders
}

func doubleDigitStr(id int) string {
  result := strconv.Itoa(id)
  if len(result) == 1 {
    result = "0" + result
  }
  return result
}
