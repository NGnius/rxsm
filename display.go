// Created 2019-07-26 by NGnius

package main

import (
	"strconv"
	"log"
	"os"

	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"
)

const (
	BUILD_MODE = 1
	PLAY_MODE = 2
)

var (
	activeSaveHandler SaveHandler
	selectedSave Save
)

// start Display
type IDisplayGoroutine interface {
	Run()
	Start()
	Join() (int, error)
}

type Display struct {
	selectedSave *Save
	activeSave *Save
	activeMode int
	activeSaves *[]Save
	saveHandler SaveHandler
	endChan chan int
	// Qt GUI objects
	window *widgets.QMainWindow
	app *widgets.QApplication
	switchModeButton *widgets.QPushButton
	currentModeLabel *widgets.QLabel
	/* TODO: import and export buttons + functionality
	importButton *widget.QPushButton2
	exportButton *widget.QPushButton2
	*/
	saveSelector *widgets.QComboBox
	newSaveButton *widgets.QPushButton
	nameField *widgets.QLineEdit
	creatorLabel *widgets.QLabel
	creatorField *widgets.QLineEdit
	// TODO: implement thumbnail button + functionality
	thumbnailImage *gui.QIcon
	thumbnailButton *widgets.QPushButton
	idLabel *widgets.QLabel
	descriptionLabel *widgets.QLabel
	descriptionField *widgets.QPlainTextEdit
	saveButton *widgets.QPushButton
	cancelButton *widgets.QPushButton
	activateButton *widgets.QPushButton
	moveButton *widgets.QPushButton
}

func NewDisplay(saveHandler SaveHandler) (*Display){
	newD := Display {endChan: make(chan int), saveHandler:saveHandler}
	newD.activeSave = saveHandler.ActiveBuildSave()
	return &newD
}

func (d *Display) Run() {
	d.activeMode = BUILD_MODE
	d.activeSaves = &d.saveHandler.BuildSaves
	log.Println(d.saveHandler.PlaySaves)
	log.Println("Display started")
	// build initial display
	d.app = widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	// with a minimum size of 250*200
	d.window = widgets.NewQMainWindow(nil, 0)
	//d.window.SetMinimumSize2(250, 200)
	d.window.SetWindowTitle("rxsm")

	d.switchModeButton = widgets.NewQPushButton2("Toggle Mode", nil)
	d.switchModeButton.ConnectClicked(d.onModeButtonClicked)
	d.currentModeLabel = widgets.NewQLabel2("Build", nil, 0)

	d.saveSelector = widgets.NewQComboBox(nil)
	d.saveSelector.AddItems(makeSelectorOptions(*d.activeSaves))
	d.saveSelector.ConnectCurrentIndexChanged(d.onSaveSelectedChanged)
	d.newSaveButton = widgets.NewQPushButton2("New", nil)
	d.newSaveButton.ConnectClicked(d.onNewSaveButtonClicked)

	d.nameField = widgets.NewQLineEdit(nil)
	d.creatorLabel = widgets.NewQLabel2("by", nil, 0)
	d.creatorField = widgets.NewQLineEdit(nil)
	d.thumbnailImage = gui.NewQIcon5("")
	d.thumbnailButton = widgets.NewQPushButton2("", d.window)
	d.thumbnailButton.SetSizePolicy(widgets.NewQSizePolicy2(4, 1, 0x00000001))
	d.thumbnailButton.ConnectClicked(d.onThumbnailButtonClicked)
	d.thumbnailButton.SetFlat(true)
	d.thumbnailButton.SetIconSize(d.thumbnailButton.Size())
	d.idLabel = widgets.NewQLabel2("id: ##", nil, 0)
	d.descriptionLabel = widgets.NewQLabel2("Description:", nil, 0)
	d.descriptionField = widgets.NewQPlainTextEdit(nil)

	d.saveButton = widgets.NewQPushButton2("Save", nil)
	d.saveButton.ConnectClicked(d.onSaveButtonClicked)
	d.cancelButton = widgets.NewQPushButton2("Cancel", nil)
	d.cancelButton.ConnectClicked(d.onCancelButtonClicked)
	d.activateButton = widgets.NewQPushButton2("Activate", nil)
	d.activateButton.ConnectClicked(d.onActivateButtonClicked)
	d.moveButton = widgets.NewQPushButton2("Toggle Location", nil)
	d.moveButton.ConnectClicked(d.onMoveToButtonClicked)

	// toggle mode to populate fields
	d.activeMode = PLAY_MODE // so toggles (back) to build mode
	d.onModeButtonClicked(true)

	headerLayout := widgets.NewQGridLayout2()
	headerLayout.AddWidget3(d.switchModeButton, 0, 0, 1, 5, 0)
	headerLayout.AddWidget3(d.currentModeLabel, 0, 5, 1, 1, 0)
	headerLayout.AddWidget3(d.saveSelector, 1, 0, 1, 5, 0)
	headerLayout.AddWidget3(d.newSaveButton, 1, 5, 1, 1, 0)

	infoLayout := widgets.NewQGridLayout2()
	infoLayout.AddWidget3(d.nameField, 0, 0, 1, 4, 0)
	infoLayout.AddWidget3(d.thumbnailButton, 0, 4, 2, 2, 0)
	infoLayout.AddWidget2(d.creatorLabel, 1, 0, 0)
	infoLayout.AddWidget3(d.creatorField, 1, 1, 1, 3, 0)

	descriptionLayout := widgets.NewQGridLayout2()
	descriptionLayout.AddWidget3(d.descriptionLabel, 0, 0, 1, 5, 0)
	descriptionLayout.AddWidget2(d.idLabel, 0, 5, 0)
	descriptionLayout.AddWidget3(d.descriptionField, 1, 0, 1, 6, 0)

	bottomButtons := widgets.NewQGridLayout2()
	bottomButtons.AddWidget2(d.saveButton, 0, 0, 0)
	bottomButtons.AddWidget2(d.cancelButton, 0, 1, 0)
	bottomButtons.AddWidget2(d.activateButton, 1, 0, 0)
	bottomButtons.AddWidget2(d.moveButton, 1, 1, 0)

	masterLayout := widgets.NewQGridLayout2()
	masterLayout.AddLayout(headerLayout, 0, 0, 0)
	masterLayout.AddLayout(infoLayout, 1, 0, 0)
	masterLayout.AddLayout(descriptionLayout, 2, 0, 0)
	masterLayout.AddLayout(bottomButtons, 3, 0, 0)

	centralWidget := widgets.NewQWidget(d.window, 0)
	centralWidget.SetLayout(masterLayout)
	d.window.SetCentralWidget(centralWidget)

	d.window.Show()
	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	d.app.Exec()
	log.Println("Display ended")
	d.endChan <- 0
}

func (d *Display) Start() {
	go d.Run()
}

func (d *Display) Join() (int, error) {
	return <- d.endChan, nil
}

func (d *Display) populateFields() {
	if d.selectedSave == nil {
		return
	}
	d.nameField.SetText(d.selectedSave.Data.Name)
	d.creatorField.SetText(d.selectedSave.Data.Creator)
	oldIdText := d.idLabel.Text()
	d.idLabel.SetText(oldIdText[:len(oldIdText)-2]+DoubleDigitStr(d.selectedSave.Data.Id))
	d.descriptionField.SetPlainText(d.selectedSave.Data.Description)
	d.thumbnailImage.Swap(gui.NewQIcon5(d.selectedSave.ThumbnailPath))
	d.thumbnailButton.SetIcon(d.thumbnailImage)
}

func (d *Display) syncBackFields() {
	d.selectedSave.Data.Name = d.nameField.Text()
	d.selectedSave.Data.Creator = d.creatorField.Text()
	d.selectedSave.Data.Description = d.descriptionField.ToPlainText()
}

func (d *Display) onModeButtonClicked(bool) {
	switch d.activeMode {
	case PLAY_MODE:
		d.currentModeLabel.SetText("build")
		d.activeMode = BUILD_MODE
		d.activeSaves = &d.saveHandler.BuildSaves
	case BUILD_MODE:
		d.activeMode = PLAY_MODE
		d.activeSaves = &d.saveHandler.PlaySaves
		d.currentModeLabel.SetText("play")
	}
	d.saveSelector.Clear()
	d.saveSelector.AddItems(makeSelectorOptions(*d.activeSaves))
	// propagation calls d.onSaveSelectedChanged(d.saveSelector.CurrentIndex())
	log.Println("Switched to mode "+strconv.Itoa(d.activeMode))
}

func (d *Display) onSaveSelectedChanged(index int) {
	if index == -1 { // no items in dropdown
		return
	}
	d.selectedSave = &(*d.activeSaves)[index]
	d.populateFields()
	log.Println("Selected "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onNewSaveButtonClicked(bool) {
	newId := d.saveHandler.MaxId() + 1
	newSave, newSaveErr := NewNewSave(d.saveHandler.PlaySaveFolderPath(newId), newId)
	if newSaveErr != nil {
		log.Println("Error while creating new save")
		log.Println(newSaveErr)
		return
	}
	switch d.activeMode {
	case BUILD_MODE:
		d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves, newSave)
	case PLAY_MODE:
		d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves, newSave)
	}
	d.saveSelector.AddItems([]string{newSave.Data.Name})
	log.Println("Created new save "+strconv.Itoa(newSave.Data.Id))
	// select newly created save
	d.saveSelector.SetCurrentIndex(len(*d.activeSaves)-1)
	// propagation calls d.onSaveSelectedChanged(len(*d.activeSaves)-1)
}

func (d *Display) onThumbnailButtonClicked(bool) {
	// TODO: implement thumbnail picker dialogue
	log.Println("Thumbnail button clicked (unimplemented)")
}

func (d *Display) onSaveButtonClicked(bool) {
	d.syncBackFields()
	saveErr := d.selectedSave.Data.Save()
	if saveErr != nil {
		log.Println(saveErr)
	}
	index := d.saveSelector.CurrentIndex()
	d.saveSelector.SetItemText(index, makeSelectorOptions(*d.activeSaves)[index])
	log.Println("Saved "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onCancelButtonClicked(bool) {
	d.populateFields()
	log.Println("Canceled "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onActivateButtonClicked(bool) {
	if d.activeMode == PLAY_MODE {
		return // button is inactive in play mode
	}
	d.activeSave.MoveToId()
	d.selectedSave.MoveToFirst()
	d.activeSave = d.selectedSave
	log.Println("Activated "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onMoveToButtonClicked(bool) {
	// TODO: implement move to opposite build/play game mode folder
	if d.selectedSave == nil {
		return
	}
	log.Println("Moving save "+strconv.Itoa(d.selectedSave.Data.Id))
	if d.selectedSave == d.activeSave {
		d.activeSave = nil
		d.selectedSave.MoveToId()
	}
	switch d.activeMode {
	case BUILD_MODE:
		moveErr := d.selectedSave.Move(d.saveHandler.PlaySaveFolderPath(d.selectedSave.Data.Id))
		if moveErr != nil {
			log.Println("Error while moving "+strconv.Itoa(d.selectedSave.Data.Id))
			log.Println(moveErr)
			return
		}
		selIndex := d.saveSelector.CurrentIndex()
		d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves, d.saveHandler.BuildSaves[selIndex]) // add to playsaves
		d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves[:selIndex], d.saveHandler.BuildSaves[selIndex+1:]...) // remove from buildsaves
	case PLAY_MODE:
		moveErr := d.selectedSave.Move(d.saveHandler.BuildSaveFolderPath(d.selectedSave.Data.Id))
		if moveErr != nil {
			log.Println("Error while moving "+strconv.Itoa(d.selectedSave.Data.Id))
			log.Println(moveErr)
			return
		}
		selIndex := d.saveSelector.CurrentIndex()
		d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves, d.saveHandler.PlaySaves[selIndex]) // add to playsaves
		d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves[:selIndex], d.saveHandler.PlaySaves[selIndex+1:]...) // remove from buildsaves
	}
	if d.activeSave == nil && len(*d.activeSaves) > 0 {
		d.activeSave = &(*d.activeSaves)[0]
		d.activeSave.MoveToFirst()
	}
	d.onModeButtonClicked(true) // toggle to other mode to keep showing selected save
	// re-select save
	d.saveSelector.SetCurrentIndex(len(*d.activeSaves)-1)
	//d.onSaveSelectedChanged(len(*d.activeSaves)-1)
	log.Println("Save moved to "+strconv.Itoa(d.activeMode))
}

// end Display

func makeSelectorOptions(saves []Save) ([]string) {
	var result []string
	for _, s := range saves {
		result = append(result, s.Data.Name)
	}
	return result
}

func moveSaveToFirst(selected *Save, saves []Save) {
	for _, s := range saves {
		err := s.MoveToId()
		if err != nil {
			log.Println(err)
		}
	}
	err := selected.MoveToFirst()
	if err != nil {
		log.Println(err)
	}
}

func getSelectedSave(name string) (save Save, isFound bool){
	var noResult Save
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
