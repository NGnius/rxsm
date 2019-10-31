// Created 2019-07-26 by NGnius

package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	RXSMVersion        string = "v2.1.0"
	RXSMPlatformStream string = "release"
	UpdateSteps        int    = 4
	DownloadTempFile          = "rxsm-update.zip"
)

var activeDisplay IDisplayGoroutine

// global update variables
var (
	IsOutOfDate bool   = false
	DownloadURL string = ""
	IsUpdating  bool   = false
)

func init() {
	log.Println("Starting init")
	// runtime.GOMAXPROCS(1)
	// load config file
	LoadGlobalConfig()
	f, logCreateErr := os.Create(GlobalConfig.LogPath)
	if logCreateErr != nil {
		log.Println("Error creating log file, skipping log.SetOutput(logFile)")
		log.Println(logCreateErr)
	} else {
		log.Println("Log directed to " + GlobalConfig.LogPath)
		log.SetOutput(f)
	}
	// log details important for debugging
	log.Println("Info for support purposes (just in case)")
	log.Println("RXSM version '" + GlobalConfig.Version + "'")
	log.Println("RXSM old version '" + GlobalConfig.LastVersion() + "'")
	log.Println("Build OS-Arch " + runtime.GOOS + "-" + runtime.GOARCH)
	log.Println("Compiler " + runtime.Compiler)
	log.Println("Processors " + strconv.Itoa(runtime.NumCPU()))
	log.Println("Init complete")
}

func main() {
	var exitVal int
	shouldExit := parseRunArgs()
	if shouldExit {
		os.Exit(exitVal)
	}
	log.Println("Starting main routine")
	GlobalConfig.Save()
	log.Println("RobocraftX Play Path: " + GlobalConfig.PlayPath)
	log.Println("RobocraftX Build Path: " + GlobalConfig.BuildPath)
	saveHandler := NewSaveHandler(GlobalConfig.PlayPath, GlobalConfig.BuildPath)
	activeDisplay = NewDisplay(saveHandler)
	activeDisplay.Start()
	exitVal, _ = activeDisplay.Join()
	if exitVal == 20 { // set new install dir
		log.Println("Display requested an update to PlayPath")
		GlobalConfig.PlayPath = filepath.FromSlash(NewInstallPath + "/" + ConfigPlayPathEnding)
		log.Println("New RobocraftX Play Path: " + GlobalConfig.PlayPath)
		GlobalConfig.Save()
		exitVal = 0
	}
	if IsUpdating {
		process, forkErr := installRXSMUpdate()
		if forkErr != nil {
			log.Println("Install failed")
			log.Println(forkErr)
		} else {
			log.Println("Forked install binary to pid", process.Pid)
		}
	}
	log.Println("rxsm terminated")
	os.Exit(exitVal) // this prevents defered operations, which may cause issues
}

func parseRunArgs() (exit bool) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-version", "--version", "version":
			versionStr := "RXSM version "
			if GlobalConfig.LastVersion() != GlobalConfig.Version {
				if GlobalConfig.LastVersion() == "" {
					versionStr += "unknown version -> "
				} else {
					versionStr += GlobalConfig.LastVersion() + " -> "
				}
			}
			versionStr += GlobalConfig.Version
			log.Println(versionStr)
			exit = true
		}
	}
	return
}

// update functions

func PlatformString() string {
	if RXSMPlatformStream == "" {
		return runtime.GOOS + "/" + runtime.GOARCH
	}
	return runtime.GOOS + "/" + runtime.GOARCH + "/" + RXSMPlatformStream
}

func checkForRXSMUpdate() (downloadURL string, isUpdateAvailable bool, ok bool) {
	downloadURL, isUpdateAvailable, ok = CheckForUpdate(GlobalConfig.UpdateServer, GlobalConfig.Version, PlatformString())
	if !ok {
		return
	}
	IsOutOfDate = isUpdateAvailable
	DownloadURL = downloadURL
	return
}

func downloadRXSMUpdateQuiet() {
	downloadRXSMUpdate(func(int, string) {})
}

func downloadRXSMUpdate(statusCallback func(progress int, description string)) {
	statusCallback(1, "Downloading")
	log.Println("Downloading update from " + DownloadURL)
	f, createErr := os.Create(DownloadTempFile)
	if createErr != nil {
		log.Println("Error creating temporary update file")
		log.Println(createErr)
		statusCallback(-1, "Error creating update temporary file")
		return
	}
	defer f.Close()
	ok := DownloadUpdate(DownloadURL, f)
	if !ok {
		log.Println("Error downloading update")
		statusCallback(-1, "Download failed")
		return
	}
	statusCallback(2, "Installing Updater")
	f.Sync()
	f.Seek(0, 0)
	fStat, statErr := f.Stat()
	if statErr != nil {
		log.Println("Error getting download temp file stat")
		log.Println(statErr)
		statusCallback(-1, "Installing failed")
		return
	}
	zipFile, zipErr := zip.NewReader(f, fStat.Size())
	if zipErr != nil {
		log.Println("Error creating zip reader")
		log.Println(zipErr)
		statusCallback(-1, "Installing failed")
		return
	}
	for _, f := range zipFile.File {
		if !f.FileHeader.Mode().IsDir() {
			filename := filepath.Clean(f.FileHeader.Name)
			var updaterFilename string
			if runtime.GOOS == "windows" {
				updaterFilename = "rxsm-updater.exe"
			} else {
				updaterFilename = "rxsm-updater"
			}
			if len(filename) >= len(updaterFilename) && filename[:len(updaterFilename)] == updaterFilename {
				statusCallback(3, "Extracting Updater")
				fileReadCloser, openErr := f.Open()
				if openErr != nil {
					log.Println("Error opening updater in zip archive")
					log.Println(openErr)
					statusCallback(-1, "Extracting failed")
					return
				}
				defer fileReadCloser.Close()
				destFile, createErr := os.Create(updaterFilename)
				if createErr != nil {
					log.Println("Error creating updater file")
					log.Println(createErr)
					statusCallback(-1, "Extracting failed")
					return
				}
				defer destFile.Close()
				_, copyErr := io.Copy(destFile, fileReadCloser)
				if copyErr != nil {
					log.Println("Error copying/extracting updater")
					log.Println(copyErr)
					statusCallback(-1, "Extracting failed")
					return
				}
			}
		}
	}
	statusCallback(UpdateSteps, "Complete")
	IsUpdating = true
}

func installRXSMUpdate() (process *os.Process, err error) {
	if runtime.GOOS == "windows" {
		return os.StartProcess(".\\rxsm-updater.exe", []string{".\\rxsm-updater.exe", "--wait", "1s", "--log", "--zip", DownloadTempFile}, nil)
	} else {
		return os.StartProcess("./rxsm-updater", []string{"./rxsm-updater", "--wait", "1s", "--log", "--zip", DownloadTempFile}, nil)
	}
}
