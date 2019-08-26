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
	BUILD_MODE = 0
	PLAY_MODE = 1
)

var (
	NewInstallPath string
	IconPath string = "icon.svg"
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
	exitStatus int
	firstTime bool
	temporaryThumbnailPath string
	endChan chan int
	// Qt GUI objects
	window *widgets.QMainWindow
	app *widgets.QApplication
	modeTab *widgets.QTabBar
	/* TODO: import and export buttons + functionality
	importButton *widget.QPushButton2
	exportButton *widget.QPushButton2
	*/
	saveSelector *widgets.QComboBox
	newSaveButton *widgets.QPushButton
	nameField *widgets.QLineEdit
	creatorLabel *widgets.QLabel
	creatorField *widgets.QLineEdit
	thumbnailImage *gui.QIcon
	thumbnailButton *widgets.QPushButton
	idLabel *widgets.QLabel
	descriptionLabel *widgets.QLabel
	descriptionField *widgets.QPlainTextEdit
	saveButton *widgets.QPushButton
	cancelButton *widgets.QPushButton
	activateCheckbox *widgets.QCheckBox
	moveButton *widgets.QPushButton
	installPathDialog *InstallPathDialog
}

func NewDisplay(saveHandler SaveHandler) (*Display){
	newD := Display {endChan: make(chan int, 1), saveHandler:saveHandler, firstTime:true}
	newD.activeSave = saveHandler.ActiveBuildSave()
	return &newD
}

func (d *Display) Run() {
	d.activeMode = BUILD_MODE
	d.activeSaves = &d.saveHandler.BuildSaves
	if d.activeSave != nil {
		log.Println("Active save on startup "+strconv.Itoa(d.activeSave.Data.Id))
	} else {
		log.Println("No active save detected on startup")
	}
	log.Println("Display started")
	// build initial display
	d.app = widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	// with a minimum size of 250*200
	d.window = widgets.NewQMainWindow(nil, 0)
	//d.window.SetMinimumSize2(250, 200)
	d.window.SetWindowTitle("rxsm")
	d.modeTab = widgets.NewQTabBar(nil)
	d.modeTab.ConnectCurrentChanged(d.onModeTabChanged)

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
	d.thumbnailButton.SetIconSize(d.thumbnailButton.Size())
	d.idLabel = widgets.NewQLabel2("ID: ##", nil, 0)
	d.descriptionLabel = widgets.NewQLabel2("Description", nil, 0)
	d.descriptionField = widgets.NewQPlainTextEdit(nil)

	d.saveButton = widgets.NewQPushButton2("Save", nil)
	d.saveButton.ConnectClicked(d.onSaveButtonClicked)
	d.cancelButton = widgets.NewQPushButton2("Cancel", nil)
	d.cancelButton.ConnectClicked(d.onCancelButtonClicked)
	d.activateCheckbox = widgets.NewQCheckBox2("Activated", nil)
	d.activateCheckbox.ConnectStateChanged(d.onActivateChecked)
	d.moveButton = widgets.NewQPushButton2("Toggle Location", nil)
	d.moveButton.ConnectClicked(d.onMoveToButtonClicked)

	// populate fields
	// propogation of events calls d.onModeTabChanged(BUILD_MODE)
	d.modeTab.AddTab("Build")
	d.modeTab.AddTab("Play")

	headerLayout := widgets.NewQGridLayout2()
	headerLayout.AddWidget3(d.modeTab, 0, 0, 1, 6, 0)
	headerLayout.AddWidget3(d.saveSelector, 1, 0, 1, 5, 0)
	headerLayout.AddWidget3(d.newSaveButton, 1, 5, 1, 1, 0)

	infoLayout := widgets.NewQGridLayout2()
	infoLayout.AddWidget3(d.activateCheckbox, 0, 0, 1, 3, 0x0004)
	infoLayout.AddWidget3(d.idLabel, 0, 3, 1, 2, 0)
	infoLayout.AddWidget3(d.nameField, 1, 0, 1, 4, 0)
	infoLayout.AddWidget3(d.thumbnailButton, 0, 4, 3, 2, 0)
	infoLayout.AddWidget2(d.creatorLabel, 2, 0, 0x0004)
	infoLayout.AddWidget3(d.creatorField, 2, 1, 1, 3, 0)

	descriptionLayout := widgets.NewQGridLayout2()
	descriptionLayout.AddWidget2(d.descriptionLabel, 0, 0, 0x0004)
	descriptionLayout.AddWidget2(d.descriptionField, 1, 0, 0)

	bottomButtons := widgets.NewQGridLayout2()
	bottomButtons.AddWidget2(d.saveButton, 0, 0, 0)
	bottomButtons.AddWidget2(d.cancelButton, 0, 1, 0)
	bottomButtons.AddWidget3(d.moveButton, 1, 0, 1, 2, 0)

	masterLayout := widgets.NewQGridLayout2()
	masterLayout.AddLayout(headerLayout, 0, 0, 0)
	masterLayout.AddLayout(infoLayout, 1, 0, 0)
	masterLayout.AddLayout(descriptionLayout, 2, 0, 0)
	masterLayout.AddLayout(bottomButtons, 3, 0, 0)

	centralWidget := widgets.NewQWidget(d.window, 0)
	centralWidget.SetLayout(masterLayout)
	d.window.SetCentralWidget(centralWidget)

	// TODO: make rxsm logo icon
	rxsmIcon := gui.NewQIcon5(IconPath)
	d.app.SetWindowIcon(rxsmIcon)

	d.window.Show()
	if len(d.saveHandler.PlaySaves) == 0 { // automatically prompt for RCX location if not default
		log.Println("PlaySaves is empty, opening install location dialog")
		d.installPathDialog = NewInstallPathDialog(d.window, 0)
		d.installPathDialog.ConnectFinished(d.onInstallPathDialogFinished)
		d.installPathDialog.OpenInstallPathDialog()
	}

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	d.app.Exec()
	d.window.Hide()
	log.Println("Display ended")
	d.endChan <- d.exitStatus
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
	d.idLabel.SetText("ID: "+DoubleDigitStr(d.selectedSave.Data.Id))
	d.descriptionField.SetPlainText(d.selectedSave.Data.Description)
	d.thumbnailImage.Swap(gui.NewQIcon5(d.selectedSave.ThumbnailPath))
	d.thumbnailButton.SetIcon(d.thumbnailImage)
	if d.activeSave == d.selectedSave {
		d.activateCheckbox.SetCheckState(2)
	} else {
		d.activateCheckbox.SetCheckState(0)
	}
}

func (d *Display) syncBackFields() {
	d.selectedSave.Data.Name = d.nameField.Text()
	d.selectedSave.Data.Creator = d.creatorField.Text()
	d.selectedSave.Data.Description = d.descriptionField.ToPlainText()
	// copy thumbnail to save location
	if d.temporaryThumbnailPath != d.selectedSave.ThumbnailPath && d.temporaryThumbnailPath != "" {
		copyErr := d.saveHandler.CopyTo(d.temporaryThumbnailPath, d.selectedSave.ThumbnailPath)
		if copyErr != nil {
			log.Print("Error while copying Thumbnail from "+d.temporaryThumbnailPath+" to "+d.selectedSave.ThumbnailPath)
			log.Println(copyErr)
		}
		d.temporaryThumbnailPath = ""
	}
}

func (d *Display) onInstallPathDialogFinished(int) {
	log.Println("Install Path Dialog closed")
	NewInstallPath = d.installPathDialog.InstallPath
	if NewInstallPath != "" {
		if d.installPathDialog.Result() == 1 {
			log.Println("New Install Path provided, requesting restart")
			d.exitStatus = 20
			// d.app.Exit(0)
		} else {
			NewInstallPath = ""
		}
	}
}

func (d *Display) onModeTabChanged(tabIndex int) {
	if tabIndex == -1 { // no tabs
		return
	}
	switch tabIndex {
	case BUILD_MODE:
		d.activeSaves = &d.saveHandler.BuildSaves
		d.activateCheckbox.SetCheckable(true)
	case PLAY_MODE:
		d.activeSaves = &d.saveHandler.PlaySaves
		d.activateCheckbox.SetCheckable(false)
	}
	d.activeMode = tabIndex
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
	var newSave Save
	var newSaveErr error
	switch d.activeMode {
	case BUILD_MODE:
		newSave, newSaveErr = NewNewSave(d.saveHandler.BuildSaveFolderPath(newId), newId)
		if newSaveErr != nil {
			log.Println("Error while creating new build save")
			log.Println(newSaveErr)
			return
		}
		d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves, newSave)
	case PLAY_MODE:
		newSave, newSaveErr = NewNewSave(d.saveHandler.PlaySaveFolderPath(newId), newId)
		if newSaveErr != nil {
			log.Println("Error while creating new play save")
			log.Println(newSaveErr)
			return
		}
		d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves, newSave)
	}
	d.saveSelector.AddItems([]string{newSave.Data.Name})
	log.Println("Created new save "+strconv.Itoa(newSave.Data.Id))
	// select newly created save
	d.saveSelector.SetCurrentIndex(len(*d.activeSaves)-1)
	// propagation calls d.onSaveSelectedChanged(len(*d.activeSaves)-1)
}

func (d *Display) onThumbnailButtonClicked(bool) {
	var fileDialog *widgets.QFileDialog = widgets.NewQFileDialog(nil, 0)
	d.temporaryThumbnailPath = fileDialog.GetOpenFileName(d.window, "Select a new thumbnail", d.selectedSave.ThumbnailPath, "Images (*.jpg)", "", 0)
	if d.temporaryThumbnailPath != "" {
		d.thumbnailImage.Swap(gui.NewQIcon5(d.temporaryThumbnailPath))
		d.thumbnailButton.SetIcon(d.thumbnailImage)
		log.Println("Thumbnail temporarily set to "+d.temporaryThumbnailPath)
	} else {
		log.Println("Thumbnail button clicked but dialog cancelled")
	}
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

func (d *Display) onActivateChecked(checkState int) {
	// TODO: implement check behaviour
	if d.activeMode == PLAY_MODE {
		log.Println("Hey buddy, you can't activate a play save")
		return // activating a play save is pointless
	}
	if d.activeSave != nil {
		moveErr := d.activeSave.MoveToId()
		if moveErr != nil {
			log.Println("Error while deactivating "+strconv.Itoa(d.activeSave.Data.Id))
			log.Println(moveErr)
			return
		}
		log.Println("Deactivated "+strconv.Itoa(d.activeSave.Data.Id))
	}
	switch checkState {
	case 0:
		d.activeSave = nil
	case 2:
		if d.selectedSave != nil {
			d.activeSave = d.selectedSave
			moveErr := d.activeSave.MoveToFirst()
			if moveErr != nil {
				log.Println("Error while activating "+strconv.Itoa(d.activeSave.Data.Id))
				log.Println(moveErr)
				return
			}
			log.Println("Activated "+strconv.Itoa(d.activeSave.Data.Id))
		} else {
			log.Println("Selected save is nil, activation failed")
		}
	}
}

func (d *Display) onMoveToButtonClicked(bool) {
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
		d.modeTab.SetCurrentIndex(PLAY_MODE) // toggle to other mode to keep showing selected save
		// propagation calls d.onModeTabChanged(PLAY_MODE)
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
		d.modeTab.SetCurrentIndex(BUILD_MODE) // toggle to other mode to keep showing selected save
		// propagation calls d.onModeTabChanged(BUILD_MODE)
	}
	if d.activeSave == nil && len(*d.activeSaves) > 0 {
		d.activeSave = &(*d.activeSaves)[0]
		d.activeSave.MoveToFirst()
	}
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
