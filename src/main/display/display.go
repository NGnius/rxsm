// Created 2019-07-26 by NGnius

package display

import (
	"strconv"
	"log"
	"strings"
	"os"
	"bufio"

  "../saver"
)

var (
	activeSaveHandler saver.SaveHandler
	selectedSave saver.Save
)

// start Display
type IDisplayGoroutine interface {
	Run()
	Start()
	Join() (int, error)
}

type Display struct {
	selectedSave saver.Save
	activeSave saver.Save
	saveHandler saver.SaveHandler
	endChan chan int
	tempText string
	tempedText string
	prependText string
	prependedText string
	currentText string
}

func NewDisplay(saveHandler saver.SaveHandler) (*Display){
	log.Println(saveHandler.PlaySaves)
	newD := Display {endChan: make(chan int), saveHandler:saveHandler}
	// set to invalid Id
	newD.selectedSave.Data.Id = -1
	newD.activeSave = saveHandler.ActiveBuildSave()
	return &newD
}

func (d *Display) Run() {
	log.Println("Display started")
	// build initial display
	d.prependText = "rxsm (RobocraftX Save Manager) -- pre-alpha v0.0.1\nTest/Dev command line version\nCommands: select #, activate, new, exit\n"
	newText := makeSelectorDisplayString(makeSelectorOptions(d.saveHandler.BuildSaves), d.selectedSave.Data.Name, d.activeSave.Data.Name)
	d.overwrite(newText)
	cliReader := bufio.NewReader(os.Stdin)
	log.Println("Entering Display input loop")
	inputLoop: for {
		text, _ := cliReader.ReadString('\n')
		text = text[:len(text)-1]
		text = strings.Trim(text, " \n\t\r")
		args := strings.Split(text, " ")
		log.Println("stdin: '"+text+"'")
		log.Println(len(args))
		switch args[0] {
		case "exit", "end":
			break inputLoop
		case "select":
			if len(args) > 1 {
				i, convErr := strconv.Atoi(args[1])
				if convErr == nil {
					i = i-1 // zero-indexed, but displayed as one-indexed
					if i < len(d.saveHandler.BuildSaves) && i > -1 {
						d.selectedSave = d.saveHandler.BuildSaves[i]
						log.Println(strconv.Itoa(i)+" selected")
					} else {
						d.tempText = "? Invalid"
					}
				}
			}
		case "activate":
			log.Println("Activating id: "+strconv.Itoa(d.selectedSave.Data.Id))
			if d.selectedSave.Data.Id == -1 {
				d.tempText = "? Please select a save first"
			} else {
				moveSaveToFirst(d.selectedSave, d.saveHandler.BuildSaves)
				d.activeSave = d.selectedSave
			}
		case "new":
			log.Println("Creating new build save")
			newId := d.saveHandler.MaxId()+1
			savePath := d.saveHandler.BuildSaveFolderPath(newId)
			log.Println("Creating new save in "+savePath)
			newSave, newSaveErr := saver.NewNewSave(savePath, newId)
			if newSaveErr != nil {
				log.Println("Error during 'new' command")
				log.Println(newSaveErr)
			}
			d.selectedSave = newSave
			d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves, d.selectedSave)
		case "help", "?":
			d.tempText = "Read the line right below this one ya big doofus :P"
		}
		d.overwrite(makeSelectorDisplayString(makeSelectorOptions(d.saveHandler.BuildSaves), d.selectedSave.Data.Name, d.activeSave.Data.Name))
	}
	log.Println("Display ended")
	d.endChan <- 0
}

func (d *Display) Start() {
	go d.Run()
}

func (d *Display) Join() (int, error) {
	return <- d.endChan, nil
}

func (d *Display) write(s string) (err error) {
	d.currentText += s
	_, err = os.Stdout.Write([]byte(s))
	return
}

func (d *Display) overwrite(s string) (err error) {
	d.clear(99, 99)
	d.currentText = ""
	if d.tempText != "" {
		err = d.write(d.tempText+"\n")
		if err != nil {
			return
		}
		d.tempedText = d.tempText+"\n"
		d.tempText = ""
	} else {
		d.tempedText = ""
	}
	err = d.write(d.prependText)
	if err != nil {
		return
	}
	d.prependedText = d.prependText
	err = d.write(s)
	return
}

func (d *Display) refresh() (err error) {
	err = d.overwrite(d.currentText[len(d.prependedText)+len(d.tempedText):])
	return
}

func (d *Display) clear(lines int, chars int) (err error) {
	d.write("\033[0;0H")
	//return
	i:=0
	loop: for {
		if i == lines {
			break loop
		}
		j:=0
		internalLoop: for {
			if j == chars {
				break internalLoop
			}
			d.write(" ")
			j++
		}
		d.write("\n")
		i++
	}
	d.write("\033[0;0H")
	d.currentText = ""
	return
}

// end Display

func makeSelectorOptions(saves []saver.Save) ([]string) {
	var result []string
	for _, s := range saves {
		result = append(result, s.Data.Name)
	}
	return result
}

func makeSelectorDisplayString(options []string, selected string, active string) (result string) {
	for i, opt := range options {
		result += strconv.Itoa(i+1)+". "+opt // zero-indexed, but displayed as one-indexed
		if opt == selected {
			result += " (selected)"
		}
		if opt == active {
			result += " (active)"
		}
		result += "\n"
	}
	return
}

func moveSaveToFirst(selected saver.Save, saves []saver.Save) {
	for _, s := range saves {
		err := s.MoveOut()
		if err != nil {
			log.Println(err)
		}
	}
	selected.MoveToFirst()
}

func getSelectedSave(name string) (save saver.Save, isFound bool){
	var noResult saver.Save
	noResult.Data.Id = -1
	for _, s := range activeSaveHandler.BuildSaves {
		if s.Data.Name == name {
			return s, true
		}
	}
	for _, s := range activeSaveHandler.PlaySaves {
		if s.Data.Name == name {
			return s, true
		}
	}
	return noResult, false
}
