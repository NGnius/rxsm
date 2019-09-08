// Created 2019-09-03 by NGnius

package main

import (
  "log"
  "time"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
  controlEnd int = 1
  controlPause int = 2
  controlResume int = 3
)

var (
  Signer [2]string = [2]string{"RXSM (automatically)", "rxsm-auto@exmods.org"}
)

func TestGit() {
  _, err := git.PlainInit("./test-data/test-git", false)
  if err != nil {
    log.Println("Error while testing go-git")
    log.Println(err)
  }
}

type ISaveVersioner interface {
  // NOTE: this is currently bound to git versioning, but there's no
  // other reason this couldn't work with other versioning software
  Repository() *git.Repository
  Worktree() *git.Worktree
  Target() *Save
  Start(p int64) // start automatic snapshots every p nanoseconds in seperate goroutine/thread
  Exit() int // end automatic snapshots, return exit code
  Pause() // pause automatic snapshots
  Resume() // resume automatic snapshots
  StageAll() // stage all save files for commit
  PlainCommit(name string) string // commit all staged changes (do nothing if nothing staged)
}

// start SaveVersioner

type SaveVersioner struct {
  save *Save
  controlChan chan int
  endChan chan int
  repo *git.Repository
  worktree *git.Worktree
  period time.Duration
  isRunning bool
}

func NewSaveVersioner(save *Save) (sv *SaveVersioner, err error) {
  var isInit bool
  sv = &SaveVersioner{controlChan: make(chan int), endChan: make(chan int), save: save}
  sv.repo, err, isInit = openOrInitGit(save.FolderPath())
  if err != nil {
    return
  }
  sv.worktree, err = sv.repo.Worktree()
  if err != nil {
    return
  }
  if isInit {
    go func(){
      sv.StageAll()
      hash := sv.PlainCommit("Initial commit")
      log.Println("Created init commit "+hash)
      }()
  }
  return
}

func (sv *SaveVersioner) Repository() (*git.Repository) {
  return sv.repo
}

func (sv *SaveVersioner) Worktree() (*git.Worktree) {
  return sv.worktree
}

func (sv *SaveVersioner) Target() (*Save) {
  return sv.save
}

func (sv *SaveVersioner) Start(p int64) {
  if p <= 1 { // automatic versioning disabled
    return
  }
  if !sv.isRunning {
    sv.period = time.Duration(p)
    sv.isRunning = true
    go sv.Run()
    log.Println("SaveVersioner spinning up")
  } else {
    log.Fatal("Cannot Start(); SaveVersioner is already running")
  }
}

func (sv *SaveVersioner) Run() {
  exitStatus := 0
  defer func(){sv.endChan <- exitStatus}()
  isPaused := false
  ticker := time.NewTicker(sv.period)
  go sv.autoCommitNow() // inits git repo
  defer ticker.Stop()
  runLoop: for {
    select {
    case c := <-sv.controlChan:
      // process control commands
      switch c {
      case controlEnd:
        break runLoop
      case controlPause:
        isPaused = true
      case controlResume:
        isPaused = false
      }
    case <-ticker.C:
      if !isPaused {
        go sv.autoCommitNow() // this is slow
      }
    }
  }
}

func (sv *SaveVersioner) Exit() (int){
  if !sv.isRunning {
    return 0
  }
  sv.controlChan <- controlEnd
  sv.isRunning = false
  return <- sv.endChan
}

func (sv *SaveVersioner) Pause() {
  if !sv.isRunning {
    return
  }
  sv.controlChan <- controlPause
}

func (sv *SaveVersioner) Resume() {
  if !sv.isRunning {
    return
  }
  sv.controlChan <- controlResume
}

func (sv *SaveVersioner) StageAll() {
  sv.worktree.Add(".")
}

func (sv *SaveVersioner) PlainCommit(name string) (string) {
  status, _ := sv.worktree.Status()
  if !status.IsClean() { // if not all files are unmodified
    // do commit
    signature := &object.Signature{Name: Signer[0], Email: Signer[1], When: time.Now()}
    hash, err := sv.worktree.Commit(name, &git.CommitOptions{Author: signature})
    if err != nil {
      log.Println("Error in SaveVersioner commit")
      log.Println(err)
      return ""
    }
    return hash.String()
  }
  return ""
}

func (sv *SaveVersioner) autoCommitNow() {
  sv.StageAll()
  sv.PlainCommit(time.Now().Format(time.RFC3339))
}

// end SaveVersioner

func openOrInitGit(folder string) (repo *git.Repository, err error, isInit bool){
  repo, err = git.PlainOpen(folder)
  if err != nil {
    isInit = true
    repo, err = git.PlainInit(folder, false)
    return
  }
  return
}
