// Created 2019-07-26 by NGnius

package main

import (
  "os"
  "path/filepath"
  "log"
  "runtime"
  "strconv"
  "fmt"
)

const (
  RXSMVersion string = "v1.0.0"
)

var activeDisplay IDisplayGoroutine

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
    log.Println("Log directed to "+GlobalConfig.LogPath)
    log.SetOutput(f)
  }
  // log details important for debugging
  log.Println("Info for support purposes (just in case)")
  log.Println("RXSM version '"+GlobalConfig.Version+"'")
  log.Println("RXSM old version '"+GlobalConfig.LastVersion()+"'")
  log.Println("Build OS-Arch "+runtime.GOOS+"-"+runtime.GOARCH)
  log.Println("Compiler "+runtime.Compiler)
  log.Println("Processors "+strconv.Itoa(runtime.NumCPU()))
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
  log.Println("RobocraftX Play Path: "+GlobalConfig.PlayPath)
  log.Println("RobocraftX Build Path: "+GlobalConfig.BuildPath)
  saveHandler := NewSaveHandler(GlobalConfig.PlayPath, GlobalConfig.BuildPath)
  activeDisplay = NewDisplay(saveHandler)
  activeDisplay.Start()
  exitVal, _ = activeDisplay.Join()
  if exitVal == 20 { // set new install dir
    log.Println("Display requested an update to PlayPath")
    GlobalConfig.PlayPath = filepath.FromSlash(NewInstallPath+"/"+ConfigPlayPathEnding)
    log.Println("New RobocraftX Play Path: "+GlobalConfig.PlayPath)
    GlobalConfig.Save()
    exitVal = 0
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
          versionStr += GlobalConfig.LastVersion()+" -> "
        }
      }
      versionStr += GlobalConfig.Version
      fmt.Println(versionStr)
      exit = true
    }
  }
  return
}
