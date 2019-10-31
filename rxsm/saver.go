// Created 2019-07-26 by NGnius

package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	GameStart     = "Game_"
	GameDataFile  = "GameData.json"
	GameSaveFile  = "GameSave.GC"
	ThumbnailFile = "Thumbnail.jpg"
	FirstFolder   = "!!!Game_00"
)

var (
	ForceUniqueIds       = false
	UsedIds              = newIdTracker()
	AppendOnCopy         = true
	StringToAppendOnCopy = " (copy)"
	AppendOnNew          = true
	StringToAppendOnNew  = " (new)"
	DefaultSavePointer   *Save
)

// start of SaveHandler
type SaveHandler struct {
	playPath   string
	buildPath  string
	PlaySaves  []Save
	BuildSaves []Save
}

func NewSaveHandler(playPath string, buildPath string) SaveHandler {
	newSaveHandler := SaveHandler{
		playPath:  filepath.FromSlash(playPath),
		buildPath: filepath.FromSlash(buildPath)}
	newSaveHandler.PlaySaves = newSaveHandler.getSaves(newSaveHandler.playPath)
	newSaveHandler.BuildSaves = newSaveHandler.getSaves(newSaveHandler.buildPath)
	return newSaveHandler
}

func (sv SaveHandler) getSaves(saveFolder string) []Save {
	var saves []Save
	folders := getFoldersInFolder(saveFolder)
	for _, folder := range folders {
		s, sErr := NewSave(folder)
		if sErr == nil {
			saves = append(saves, s)
		} else {
			log.Println("Error in SaveHandler.getSaves")
			log.Println(sErr)
		}
	}
	return saves
}

func (sv SaveHandler) PlaySaveFolderPath(id int) string {
	return filepath.Join(sv.playPath, GameStart+DoubleDigitStr(id))
}

func (sv SaveHandler) BuildSaveFolderPath(id int) string {
	return filepath.Join(sv.buildPath, GameStart+DoubleDigitStr(id))
}

func (sv SaveHandler) FirstBuildSaveFolderPath() string {
	return filepath.Join(sv.buildPath, FirstFolder)
}

func (sv SaveHandler) MaxId() int {
	return UsedIds.max()
}

func (sv SaveHandler) ActiveBuildSave() (as *Save) {
	firstSaveFolder := sv.FirstBuildSaveFolderPath()
	for _, save := range sv.BuildSaves {
		if len(save.FolderPath()) >= len(firstSaveFolder) && save.FolderPath()[:len(firstSaveFolder)] == firstSaveFolder {
			as = &save
			return
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

func (sv SaveHandler) PlayPath() string {
	return sv.playPath
}

func (sv SaveHandler) BuildPath() string {
	return sv.buildPath
}

// end of SaveHandler

// start Save
type Save struct {
	Data          *GameData
	dataPath      string
	savePath      string
	thumbnailPath string
	folder        string
}

func NewSave(folder string) (Save, error) {
	newSave := Save{
		dataPath:      filepath.Join(folder, GameDataFile),
		savePath:      filepath.Join(folder, GameSaveFile),
		thumbnailPath: filepath.Join(folder, ThumbnailFile),
		folder:        folder}
	newGD, gdErr := NewGameData(newSave.dataPath)
	newSave.Data = newGD
	if gdErr != nil {
		return newSave, gdErr
	}
	// force unique ids
	if UsedIds.contains(newSave.Data.Id) {
		log.Println("Duplicate id " + strconv.Itoa(newSave.Data.Id))
		if GlobalConfig.ForceUniqueIds {
			newSave.Data.Id = UsedIds.max() + 1
			newSave.Data.Save()
		}
	}
	UsedIds.add(newSave.Data.Id)
	return newSave, nil
}

func NewNewSave(folder string, id int) (newSave Save, err error) {
	// duplicate default save
	if DefaultSavePointer == nil { // init if necessary
		var tmpSave Save
		tmpSave, err = NewSave(GlobalConfig.DefaultSaveFolder)
		if err != nil {
			return
		}
		DefaultSavePointer = &tmpSave
	}
	newSave, err = DefaultSavePointer.Duplicate(folder, id)
	if err != nil {
		return
	}
	if AppendOnNew {
		newSave.Data.Name = DefaultSavePointer.Data.Name + StringToAppendOnNew
		newSave.Data.Save()
	}
	newSave.Data.Creator = GlobalConfig.Creator
	return
}

func (s *Save) Duplicate(newFolder string, id int) (newSave Save, err error) {
	err = os.MkdirAll(newFolder, os.ModeDir|os.ModePerm)
	if err != nil {
		return
	}
	UsedIds.add(id)
	toDuplicate := [][]string{ // { {source, dest}, ...}
		{s.dataPath, filepath.Join(newFolder, GameDataFile)},
		{s.savePath, filepath.Join(newFolder, GameSaveFile)},
		{s.thumbnailPath, filepath.Join(newFolder, ThumbnailFile)}}
	for _, dupPair := range toDuplicate {
		src, openErr := os.Open(dupPair[0])
		defer src.Close()
		if openErr != nil {
			return newSave, openErr
		}
		dst, openErr := os.Create(dupPair[1])
		defer dst.Close()
		if openErr != nil {
			return newSave, openErr
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			return
		}
		dst.Sync()
	}
	// load copied save
	newSave, err = NewSave(newFolder)
	if err != nil {
		return
	}
	newSave.Data.Id = id
	if AppendOnCopy {
		newSave.Data.Name += StringToAppendOnCopy
	}
	err = newSave.Data.Save()
	return
}

func (s *Save) Move(to string) error {
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

func (s *Save) MoveToId() error {
	idDir, _ := filepath.Split(s.folder)
	idDir = filepath.Join(idDir, GameStart+DoubleDigitStr(s.Data.Id))
	if idDir == s.folder { // don't move out if already in id spot
		return nil
	}
	return s.Move(idDir)
}

func (s *Save) MoveToFirst() error {
	firstDir, _ := filepath.Split(s.folder)
	firstDir = filepath.Join(firstDir, FirstFolder)
	return s.Move(firstDir)
}

func (s *Save) MoveOut() error {
	exists := true
	rootFolder, _ := filepath.Split(s.folder)
	lastDir := filepath.Join(rootFolder, GameStart+DoubleDigitStr(s.Data.Id))
	if s.folder == lastDir && s.Data.Id != 0 { // don't move out if already in non-zero place
		return nil
	}
dirCheckLoop:
	for {
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

func (s *Save) FreeID() {
	UsedIds.remove(s.Data.Id)
}

func (s *Save) DataPath() string {
	return s.dataPath
}

func (s *Save) SavePath() string {
	return s.savePath
}

func (s *Save) ThumbnailPath() string {
	return s.thumbnailPath
}

func (s *Save) FolderPath() string {
	return s.folder
}

// end Save

// start of GameData
type GameData struct { // reference: GameData.json
	Id          int    `json:"GameID"`
	Name        string `json:"GameName"`
	Description string `json:"GameDescription"`
	Creator     string `json:"CreatorName"`
	path        string
	isInited    bool
}

func NewGameData(path string) (*GameData, error) {
	var gd GameData
	f, openErr := os.Open(path)
	if openErr != nil {
		return &gd, openErr
	}
	data, _ := ioutil.ReadAll(f)
	cleanData := cleanJsonBytes(data)
	marshErr := json.Unmarshal(cleanData, &gd)
	if marshErr != nil {
		return &gd, marshErr
	}
	// check forced values
	if GlobalConfig.ForceCreator {
		gd.Creator = GlobalConfig.Creator
	}
	gd.path = path
	gd.isInited = true
	return &gd, nil
}

func (gd *GameData) Save() error {
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

// start of idTracker

type idTracker struct {
	idMap   map[int]bool
	idArray []int
	lock    sync.Mutex
}

func newIdTracker() *idTracker {
	tracker := idTracker{idMap: map[int]bool{}}
	return &tracker
}

func (it *idTracker) add(id int) {
	log.Println("Tracker add " + strconv.Itoa(id))
	it.lock.Lock()
	it.idMap[id] = true
	it.idArray = append(it.idArray, id)
	it.lock.Unlock()
}

func (it *idTracker) remove(id int) {
	log.Println("Tracker remove " + strconv.Itoa(id))
	it.lock.Lock()
	delete(it.idMap, id)
	loc := it._location(id)
	if loc == -1 {
		log.Println("Id not found, skipping remove")
		it.lock.Unlock()
		return
	}
	if loc == 0 {
		it.idArray = it.idArray[1:]
	} else if len(it.idArray) == loc+1 {
		it.idArray = it.idArray[:loc]
	} else {
		it.idArray = append(it.idArray[:loc], it.idArray[:loc+1]...)
	}
	it.lock.Unlock()
}

func (it *idTracker) contains(id int) bool {
	it.lock.Lock()
	_, ok := it.idMap[id]
	it.lock.Unlock()
	return ok
}

func (it *idTracker) max() (max int) {
	it.lock.Lock()
	defer it.lock.Unlock()
	for _, elem := range it.idArray {
		if elem > max {
			max = elem
		}
	}
	return max
}

func (it *idTracker) _location(id int) int {
	for i, elem := range it.idArray {
		if elem == id {
			return i
		}
	}
	return -1
}

// end of idTracker

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
		if statErr == nil && info.IsDir() { // ignore errored pathes
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

func cleanJsonBytes(data []byte) (cleanData []byte) {
	var isInString bool = false
	// clean data; replace '\r' with '\\r', etc.
	double_quote := ([]byte("\""))[0]
	carriage_return := ([]byte("\r"))[0]
	page_return := ([]byte("\n"))[0]
	tab := ([]byte("\t"))[0]
	backslash := ([]byte("\\"))[0]
	for i, b := range data {
		if b == double_quote && data[i-1] != backslash {
			isInString = !isInString
		}
		if isInString {
			switch b {
			case carriage_return:
				cleanData = append(cleanData, backslash, ([]byte("r"))[0])
			case page_return:
				cleanData = append(cleanData, backslash, ([]byte("n"))[0])
			case tab:
				cleanData = append(cleanData, backslash, ([]byte("t"))[0])
			default:
				cleanData = append(cleanData, b)
			}
		} else {
			cleanData = append(cleanData, b)
		}
	}
	return
}
