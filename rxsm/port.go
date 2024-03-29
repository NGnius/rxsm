// Created 2019-08-28 by NGnius
// port as in export/import

package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// NOTE: zip requires forward slashes (/) no matter the OS
// If only Windows worked like that...

func Export(path string, save Save) error {
	// save save to a zip archive located at path
	file, createErr := os.Create(path)
	if createErr != nil {
		return createErr
	}
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()
	// create game save folder
	folderPath := GameStart + DoubleDigitStr(save.Data.Id) + "/"
	_, folderErr := zipWriter.Create(folderPath)
	if folderErr != nil {
		return folderErr
	}
	// create & write GameData.json
	data, dataReadErr := readAllFromPath(save.DataPath())
	if dataReadErr != nil {
		return dataReadErr
	}
	dataWriteErr := writeToZip(folderPath+GameDataFile, data, zipWriter)
	if dataWriteErr != nil {
		return dataWriteErr
	}
	// create & write GameSave.GC
	saveData, saveReadErr := readAllFromPath(save.SavePath())
	if saveReadErr != nil {
		return saveReadErr
	}
	saveWriteErr := writeToZip(folderPath+GameSaveFile, saveData, zipWriter)
	if saveWriteErr != nil {
		return saveWriteErr
	}
	// create & write Thumbnail.jpg
	thumbnail, thumbReadErr := readAllFromPath(save.ThumbnailPath())
	if thumbReadErr != nil {
		return thumbReadErr
	}
	thumbWriteErr := writeToZip(folderPath+ThumbnailFile, thumbnail, zipWriter)
	return thumbWriteErr
}

func Import(path string, outFolder string) (saves []*Save, err error) {
	// Load the saves contained in a zip archive located at path
	var readCloser *zip.ReadCloser
	readCloser, err = zip.OpenReader(path)
	defer readCloser.Close()
	if err != nil {
		return
	}
	candidates := map[string]map[string]*zip.File{}
	resultChan := make(chan *Save)
	for _, f := range readCloser.Reader.File {
		if !f.FileHeader.Mode().IsDir() {
			baseFolder, filename := filepath.Split(f.FileHeader.Name)
			submap, ok := candidates[baseFolder]
			if !ok {
				submap = map[string]*zip.File{}
			}
			if filename == GameDataFile || filename == GameSaveFile || filename == ThumbnailFile {
				submap[filename] = f
				ok = true
			}
			if ok {
				candidates[baseFolder] = submap
			}
		}
	}
	err = os.MkdirAll(outFolder, os.ModeDir|os.ModePerm)
	if err != nil {
		return
	}
	workers := 0
	for _, fileMap := range candidates {
		_, ok := fileMap[GameSaveFile]
		if ok {
			forcedId := UsedIds.max() + 1 + workers
			tmpFolder := filepath.Join(outFolder, GameStart+strconv.Itoa(forcedId))
			go extractSaveWorker(tmpFolder, fileMap, resultChan, forcedId)
			workers++
		}
	}
	for i := 0; i < workers; i++ {
		newSave := <-resultChan
		if newSave != nil {
			saves = append(saves, newSave)
		}
	}
	return
}

func extractSaveWorker(outFolder string, fileMap map[string]*zip.File, outChan chan *Save, forcedId int) {
	// extracts save files to correct folder
	// replace missing data with data from DefaultSaveFolder
	// extract/create GameData.json
	makeDirErr := os.Mkdir(outFolder, os.ModeDir|os.ModePerm)
	if makeDirErr != nil {
		os.RemoveAll(outFolder)
		log.Println("Failed to make extraction target directory")
		log.Println(makeDirErr)
		outChan <- nil
		return
	}
	gameDataErr := extractOrCreateFile(outFolder, fileMap, GameDataFile)
	if gameDataErr != nil {
		os.RemoveAll(outFolder)
		log.Println("GameData extraction/create err")
		log.Println(gameDataErr)
		outChan <- nil
		return
	}
	// extract GameSave.GC
	gameSaveSrc, sOpenErr := fileMap[GameSaveFile].Open() // assume exists
	defer gameSaveSrc.Close()
	if sOpenErr != nil {
		os.RemoveAll(outFolder)
		log.Println("GameSave open err")
		log.Println(sOpenErr)
		outChan <- nil
		return
	}
	gameSaveDest, sCreateErr := os.Create(filepath.Join(outFolder, GameSaveFile))
	defer gameSaveDest.Close()
	if sCreateErr != nil {
		os.RemoveAll(outFolder)
		log.Println("GameSave create err")
		log.Println(sCreateErr)
		outChan <- nil
		return
	}
	_, sCopyErr := io.Copy(gameSaveDest, gameSaveSrc)
	if sCopyErr != nil {
		os.RemoveAll(outFolder)
		log.Println("GameSave copy err")
		log.Println(sCopyErr)
		outChan <- nil
		return
	}
	// extract/create Thumbnail.jpg
	thumbnailErr := extractOrCreateFile(outFolder, fileMap, ThumbnailFile)
	if thumbnailErr != nil {
		os.RemoveAll(outFolder)
		log.Println("Thumbnail extract/create err")
		log.Println(thumbnailErr)
		outChan <- nil
		return
	}
	// create save
	newSave, newSaveErr := NewSave(outFolder)
	if newSaveErr != nil {
		os.RemoveAll(outFolder)
		log.Println("Extracted Save load err")
		log.Println(newSaveErr)
		outChan <- nil
		return
	}
	newSave.Data.Id = forcedId
	newSave.Data.Save()
	outChan <- &newSave
	return
}

func extractOrCreateFile(outFolder string, fileMap map[string]*zip.File, name string) error {
	dataZipFile, zipFileExists := fileMap[name]
	dataDest, dataDestCreateErr := os.Create(filepath.Join(outFolder, name))
	defer dataDest.Close()
	if dataDestCreateErr != nil {
		return dataDestCreateErr
	}
	var dataSrc io.ReadCloser
	var dataSrcOpenErr error
	if zipFileExists {
		dataSrc, dataSrcOpenErr = dataZipFile.Open()
	} else {
		dataSrc, dataSrcOpenErr = os.Open(filepath.Join(GlobalConfig.DefaultSaveFolder, name))
	}
	defer dataSrc.Close()
	if dataSrcOpenErr != nil {
		return dataSrcOpenErr
	}
	_, copyErr := io.Copy(dataDest, dataSrc)
	if copyErr != nil {
		return copyErr
	}
	return nil
}

func readAllFromPath(path string) (data []byte, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(file)
	return
}

func writeToZip(path string, data []byte, archiveZip *zip.Writer) (err error) {
	var fileWriter io.Writer
	fileWriter, err = archiveZip.Create(path)
	if err != nil {
		return
	}
	_, err = fileWriter.Write(data)
	return
}

func writeToPath(path string, data []byte) (err error) {
	var file *os.File
	file, err = os.Create(path)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	return
}

/*
func ioCopy(src io.Reader, dest io.Writer) (error){
  data, readErr := ioutil.ReadAll(src)
  if readErr != nil {
    return readErr
  }
  _, writeErr := dest.Write(data)
  if writeErr != nil {
    return writeErr
  }
  return nil
}*/
