// Created 2019-07-26 by NGnius

package main

import (
  "os"
  "io/ioutil"
  "io"
  "path/filepath"
  "encoding/json"
  "strconv"
  "log"
)

const (
  GameStart = "Game_"
  GameDataFile = "GameData.json"
  GameSaveFile = "GameSave.RCX"
  ThumbnailFile = "Thumbnail.jpg"
)

var (
  ForceGameCreator = false
  ForceGameCreatorTo = ""
  DefaultSaveFolder = "resources/default_save"
  TempNewSaveFolder = "tempsave"
)

// start of SaveHandler
type SaveHandler struct {
  playPath string
  buildPath string
  PlaySaves []Save
  BuildSaves []Save
}

func NewSaveHandler(playPath string, buildPath string) (SaveHandler) {
  newSaveHandler := SaveHandler {
    playPath: filepath.FromSlash(playPath),
    buildPath: filepath.FromSlash(buildPath)}
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
    } else {
      log.Println(sErr)
    }
  }
  return saves
}

func (sv SaveHandler) PlaySaveFolderPath(id int) (string) {
 return filepath.Join(sv.playPath, GameStart+DoubleDigitStr(id))
}

func (sv SaveHandler) BuildSaveFolderPath(id int) (string) {
  return filepath.Join(sv.buildPath, GameStart+DoubleDigitStr(id))
}

func (sv SaveHandler) MaxId() (max int) {
  max = -1
  for _, save := range sv.BuildSaves {
    if max < save.Data.Id {
      max = save.Data.Id
    }
  }
  for _, save := range sv.PlaySaves {
    if max < save.Data.Id {
      max = save.Data.Id
    }
  }
  return
}

func (sv SaveHandler) ActiveBuildSave() (as *Save) {
  for _, save := range sv.BuildSaves {
    if save.folder == sv.BuildSaveFolderPath(0) {
      as = &save
    }
  }
  return
}

func (sv SaveHandler) CopyTo(from string, to string) (copyErr error) {
  var fileBytes []byte
  var in, out *os.File
  in, copyErr = os.Open(from)
  if copyErr != nil {
    return
  }
  out, copyErr = os.Create(to)
  if copyErr != nil {
    return
  }
  fileBytes, copyErr = ioutil.ReadAll(in)
  if copyErr != nil {
    return
  }
  out.Write(fileBytes)
  in.Close()
  out.Sync()
  out.Close()
  return
}
// end of SaveHandler

// start Save
type Save struct {
  Data *GameData
  dataPath string
  savePath string
  ThumbnailPath string
  folder string
}

func NewSave(folder string) (Save, error) {
  newSave := Save {
    dataPath: filepath.Join(folder, GameDataFile),
    savePath: filepath.Join(folder, GameSaveFile),
    ThumbnailPath: filepath.Join(folder, ThumbnailFile),
    folder: folder}
  newGD, gdErr := NewGameData(newSave.dataPath)
  newSave.Data = newGD
  if gdErr != nil {
    return newSave, gdErr
  }
  return newSave, nil
}

func NewNewSave(folder string, id int) (newSave Save, err error) {
  // duplicate default save
  stat, statErr := os.Stat(folder)
  if statErr != nil || os.IsNotExist(statErr) {
    err = os.MkdirAll(folder, os.ModeDir|os.ModePerm)
    if err != nil {
      return
    }
  }
  if statErr == nil && !stat.IsDir() {
    log.Println("temp dir exists but is not folder, NewNewSave may fail with no error")
  }
  toDuplicate := [][]string { // { {source, dest}, ...}
    []string {filepath.Join(DefaultSaveFolder, GameDataFile), filepath.Join(folder, GameDataFile)},
    []string {filepath.Join(DefaultSaveFolder, GameSaveFile), filepath.Join(folder, GameSaveFile)},
    []string {filepath.Join(DefaultSaveFolder, ThumbnailFile), filepath.Join(folder, ThumbnailFile)}}
  for _, dupPair := range toDuplicate {
    src, openErr := os.Open(dupPair[0])
    if openErr != nil {
      return newSave, openErr
    }
    dst, openErr := os.Create(dupPair[1])
    if openErr != nil {
      return newSave, openErr
    }
    _, err = io.Copy(dst, src)
    if err != nil {
      return
    }
    dst.Sync()
    dst.Close()
  }
  // load copied save
  newSave, err = NewSave(folder)
  if err != nil {
    return
  }
  newSave.Data.Id = id
  err = newSave.Data.Save()
  return
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
  s.ThumbnailPath = filepath.Join(to, ThumbnailFile)
  return nil
}

func (s *Save) MoveToId() (error) {
  idDir, _ := filepath.Split(s.folder)
  idDir = filepath.Join(idDir, GameStart+DoubleDigitStr(s.Data.Id))
  if idDir == s.folder { // don't move out if already in id spot
    return nil
  }
  return s.Move(idDir)
}

func (s *Save) MoveToFirst() (error) {
  firstDir, _ := filepath.Split(s.folder)
  firstDir = filepath.Join(firstDir, "!!!"+GameStart+"00")
  return s.Move(firstDir)
}

func (s *Save) MoveOut() (error) {
  exists := true
  rootFolder, _ := filepath.Split(s.folder)
  lastDir := filepath.Join(rootFolder, GameStart+DoubleDigitStr(s.Data.Id))
  if s.folder == lastDir && s.Data.Id != 0 { // don't move out if already in non-zero place
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
      lastDir = filepath.Join(rootFolder, GameStart+DoubleDigitStr(s.Data.Id))
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

func NewGameData(path string) (*GameData, error) {
  var gd GameData
  f, openErr := os.Open(path)
  if openErr != nil {
    return &gd, openErr
  }
  data, _ := ioutil.ReadAll(f)
  marshErr := json.Unmarshal(data, &gd)
  if marshErr != nil {
    return &gd, marshErr
  }
  // check forced values
  if ForceGameCreator {
    gd.Creator = ForceGameCreatorTo
  }
  gd.path = path
  gd.isInited = true
  return &gd, nil
}

func (gd *GameData) Save() (error) {
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
    log.Println("File error in getFoldersInFolder")
    log.Println(dirErr)
    return folders
  }
  paths, readErr := dir.Readdirnames(0)
  if readErr != nil {
    log.Println("Read error in getFoldersInFolder")
    log.Println(readErr)
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

func DoubleDigitStr(id int) string {
  result := strconv.Itoa(id)
  if len(result) == 1 {
    result = "0" + result
  }
  return result
}
