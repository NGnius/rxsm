// Created 2019-07-26 by NGnius

package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time" // import performance stats

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	BUILD_MODE = 0
	PLAY_MODE  = 1
)

var (
	NewInstallPath   string
	SettingsIconPath = filepath.FromSlash("gear.svg")
	NewIconPath      = filepath.FromSlash("new.svg")
	ImportIconPath   = filepath.FromSlash("import-zip.svg")
	ExportIconPath   = filepath.FromSlash("export-zip.svg")
	CopyIconPath     = filepath.FromSlash("duplicate.svg")
	SaveIconPath     = filepath.FromSlash("floppy.svg")
	CancelIconPath   = filepath.FromSlash("cancel.svg")
	ActiveIconPath   = filepath.FromSlash("active.svg")
	ToggleIconPath   = filepath.FromSlash("")
	VersionsIconPath = filepath.FromSlash("fork.svg")
)

// start Display
type IDisplayGoroutine interface {
	Run()
	Start()
	Join() (int, error)
}

type Display struct {
	selectedSave           *Save
	activeSave             *Save
	activeMode             int
	activeSaves            *[]Save
	filteredSaves          []Save
	filterMapping          map[int]int
	saveHandler            SaveHandler
	saveVersioner          ISaveVersioner
	exitStatus             int
	firstTime              bool
	temporaryThumbnailPath string
	endChan                chan int
	// Qt GUI objects
	window               *widgets.QMainWindow
	app                  *widgets.QApplication
	modeTab              *widgets.QTabBar
	settingsButton       *widgets.QPushButton
	settingsIcon         *gui.QIcon
	importButton         *widgets.QPushButton
	exportButton         *widgets.QPushButton
	selectionTabWidget   *widgets.QTabWidget
	saveSelector         *widgets.QComboBox
	copySaveButton       *widgets.QPushButton
	newSaveButton        *widgets.QPushButton
	filteredSaveSelector *widgets.QComboBox
	resetFilterButton    *widgets.QPushButton
	applyFilterButton    *widgets.QPushButton
	filterTypeSelector   *widgets.QComboBox
	filterSearchField    *widgets.QLineEdit
	nameField            *widgets.QLineEdit
	creatorLabel         *widgets.QLabel
	creatorField         *widgets.QLineEdit
	thumbnailImage       *gui.QIcon
	thumbnailButton      *widgets.QPushButton
	idLabel              *widgets.QLabel
	descriptionLabel     *widgets.QLabel
	descriptionField     *widgets.QPlainTextEdit
	saveButton           *widgets.QPushButton
	cancelButton         *widgets.QPushButton
	activateCheckbox     *widgets.QCheckBox
	moveButton           *widgets.QPushButton
	versionsButton       *widgets.QPushButton
	installPathDialog    *InstallPathDialog
	settingsDialog       *SettingsDialog
}

func NewDisplay(saveHandler SaveHandler) *Display {
	newD := Display{endChan: make(chan int, 1), saveHandler: saveHandler, firstTime: true}
	newD.activeSave = saveHandler.ActiveBuildSave()
	return &newD
}

func (d *Display) Run() {
	d.activeMode = BUILD_MODE
	d.activeSaves = &d.saveHandler.BuildSaves
	if d.activeSave != nil {
		log.Println("Active save on startup " + strconv.Itoa(d.activeSave.Data.Id))
		sv, svErr := NewSaveVersioner(d.activeSave)
		if svErr != nil {
			log.Println("Error creating SaveVersioner for active save")
			log.Println(svErr)
		} else {
			d.saveVersioner = sv
			d.saveVersioner.Start(GlobalConfig.SnapshotPeriod)
		}
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
	d.settingsButton = widgets.NewQPushButton2("", nil)
	// prefer theme icon, but fallback to RXSM settings icon
	var fallBackIcon *gui.QIcon = gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, SettingsIconPath))
	d.settingsIcon = fallBackIcon.FromTheme("settings")
	if d.settingsIcon.IsNull() {
		log.Println("Falling back to RXSM settings icon")
		d.settingsIcon = fallBackIcon
	}
	d.settingsButton.SetIcon(d.settingsIcon)
	d.settingsButton.SetToolTip("Settings")
	d.settingsButton.ConnectClicked(d.onSettingsButtonClicked)
	d.importButton = widgets.NewQPushButton2("Import", nil)
	d.importButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, ImportIconPath)))
	d.importButton.SetToolTip("Import saves from a zip file")
	d.importButton.ConnectClicked(d.onImportButtonClicked)
	d.exportButton = widgets.NewQPushButton2("Export", nil)
	d.exportButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, ExportIconPath)))
	d.exportButton.SetToolTip("Export selected save to a zip file")
	d.exportButton.ConnectClicked(d.onExportButtonClicked)

	// select tab
	selectionLabel := widgets.NewQLabel2("Select a save...", nil, 0)
	// selectionLabel.SetSizePolicy2(4, 4) // shrink; don't get bigger
	d.saveSelector = widgets.NewQComboBox(nil)
	d.saveSelector.AddItems(makeSelectorOptions(*d.activeSaves))
	d.saveSelector.SetToolTip("Selected save")
	d.saveSelector.ConnectCurrentIndexChanged(d.onSaveSelectedChanged)
	d.copySaveButton = widgets.NewQPushButton2("Copy", nil)
	d.copySaveButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, ActiveIconPath)))
	d.copySaveButton.SetToolTip("Duplicate the selected save")
	d.copySaveButton.ConnectClicked(d.onCopySaveButtonClicked)
	d.newSaveButton = widgets.NewQPushButton2("New", nil)
	d.newSaveButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, NewIconPath)))
	d.newSaveButton.SetToolTip("Create a new save, using default_save")
	d.newSaveButton.ConnectClicked(d.onNewSaveButtonClicked)

	selectorMasterWidget := widgets.NewQWidget(nil, 0)
	selectionLayout := widgets.NewQGridLayout2()
	selectionLayout.AddWidget2(selectionLabel, 0, 0, 0)
	selectionLayout.AddWidget3(d.saveSelector, 1, 0, 1, 2, 0)
	selectionLayout.AddWidget3(d.newSaveButton, 2, 0, 1, 1, 0)
	selectionLayout.AddWidget3(d.copySaveButton, 2, 1, 1, 1, 0)
	selectionLayout.AddWidget3(d.importButton, 3, 0, 1, 1, 0)
	selectionLayout.AddWidget3(d.exportButton, 3, 1, 1, 1, 0)
	selectorMasterWidget.SetLayout(selectionLayout)

	// filter tab
	filterLabel := widgets.NewQLabel2("Find a save...", nil, 0)
	d.filteredSaveSelector = widgets.NewQComboBox(nil)
	d.filteredSaveSelector.SetToolTip("Filtered save results")
	d.filteredSaveSelector.ConnectCurrentIndexChanged(d.onFilterSaveSelectedChanged)
	d.applyFilterButton = widgets.NewQPushButton2("Apply", nil)
	d.applyFilterButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, ActiveIconPath)))
	d.applyFilterButton.SetToolTip("Search according to the filter")
	d.applyFilterButton.ConnectClicked(d.onApplyFilterButtonClicked)
	d.resetFilterButton = widgets.NewQPushButton2("Reset", nil)
	d.resetFilterButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, CancelIconPath)))
	d.resetFilterButton.SetToolTip("Reset search filter")
	d.resetFilterButton.ConnectClicked(d.onResetFilterButtonClicked)
	d.filterTypeSelector = widgets.NewQComboBox(nil)
	d.filterTypeSelector.AddItems([]string{"Any", "Name", "Creator", "Description", "ID"})
	d.filterTypeSelector.SetToolTip("Parameter matching method")
	d.filterSearchField = widgets.NewQLineEdit(nil)
	d.filterSearchField.SetToolTip("Parameter to search for")

	filterMasterWidget := widgets.NewQWidget(nil, 0)
	filterLayout := widgets.NewQGridLayout2()
	filterLayout.AddWidget2(filterLabel, 0, 0, 0)
	filterLayout.AddWidget3(d.filterTypeSelector, 1, 0, 1, 1, 0)
	filterLayout.AddWidget3(d.filterSearchField, 1, 1, 1, 8, 0)
	filterLayout.AddWidget3(d.applyFilterButton, 2, 0, 1, 5, 0)
	filterLayout.AddWidget3(d.resetFilterButton, 2, 5, 1, 4, 0)
	filterLayout.AddWidget3(d.filteredSaveSelector, 3, 0, 1, 9, 0)

	filterMasterWidget.SetLayout(filterLayout)

	// selection tab widget init
	d.selectionTabWidget = widgets.NewQTabWidget(nil)
	d.selectionTabWidget.SetTabPosition(2) // west
	//d.selectionTabWidget.SetSizePolicy2(1, 4) // shrink vertical
	d.selectionTabWidget.AddTab(selectorMasterWidget, "Load")
	d.selectionTabWidget.AddTab(filterMasterWidget, "Find")
	d.onResetFilterButtonClicked(false)

	d.nameField = widgets.NewQLineEdit(nil)
	d.creatorLabel = widgets.NewQLabel2("by", nil, 0)
	d.creatorField = widgets.NewQLineEdit(nil)
	d.thumbnailImage = gui.NewQIcon5("")
	d.thumbnailButton = widgets.NewQPushButton2("", d.window)
	d.thumbnailButton.SetSizePolicy(widgets.NewQSizePolicy2(2, 1, 0x00000001))
	d.thumbnailButton.SetToolTip("Select the JPG thumbnail for the selected save")
	d.thumbnailButton.ConnectClicked(d.onThumbnailButtonClicked)
	d.thumbnailButton.SetIconSize(d.thumbnailButton.Size())
	d.idLabel = widgets.NewQLabel2("ID: ##", nil, 0)
	d.idLabel.SetToolTip("GameID of the selected save")
	d.descriptionLabel = widgets.NewQLabel2("Description", nil, 0)
	d.descriptionField = widgets.NewQPlainTextEdit(nil)
	d.descriptionField.SetToolTip("Game save description")

	d.saveButton = widgets.NewQPushButton2("Save", nil)
	d.saveButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, SaveIconPath)))
	d.saveButton.SetToolTip("Save changes to input fields")
	d.saveButton.ConnectClicked(d.onSaveButtonClicked)
	d.cancelButton = widgets.NewQPushButton2("Cancel", nil)
	d.cancelButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, CancelIconPath)))
	d.cancelButton.SetToolTip("Revert changes to input fields")
	d.cancelButton.ConnectClicked(d.onCancelButtonClicked)
	d.activateCheckbox = widgets.NewQCheckBox2("Activated", nil)
	d.activateCheckbox.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, ActiveIconPath)))
	d.activateCheckbox.SetToolTip("Since Experiment 9, this has no effect")
	d.activateCheckbox.ConnectStateChanged(d.onActivateChecked)
	d.moveButton = widgets.NewQPushButton2("Toggle Location", nil)
	d.moveButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, ToggleIconPath)))
	d.moveButton.SetToolTip("Move the selected save")
	d.moveButton.ConnectClicked(d.onMoveToButtonClicked)
	d.versionsButton = widgets.NewQPushButton2("Versions", nil)
	d.versionsButton.SetIcon(gui.NewQIcon5(filepath.Join(GlobalConfig.IconPackPath, VersionsIconPath)))
	d.versionsButton.SetToolTip("View revisions of the active save")
	d.versionsButton.ConnectClicked(d.onVersionsButtonClicked)

	// populate fields
	// propogation of events calls d.onModeTabChanged(BUILD_MODE)
	d.modeTab.AddTab("Build")
	d.modeTab.AddTab("Play")

	headerLayout := widgets.NewQGridLayout2()
	headerLayout.AddWidget3(d.modeTab, 0, 0, 1, 8, 0)
	headerLayout.AddWidget2(d.settingsButton, 0, 8, 0)
	headerLayout.AddWidget3(d.selectionTabWidget, 1, 0, 1, 9, 0)

	infoLayout := widgets.NewQGridLayout2()
	infoLayout.AddWidget3(d.activateCheckbox, 0, 0, 1, 3, 0x0004)
	infoLayout.AddWidget3(d.idLabel, 0, 3, 1, 1, 0)
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
	bottomButtons.AddWidget2(d.moveButton, 1, 0, 0)
	bottomButtons.AddWidget2(d.versionsButton, 1, 1, 0)

	masterLayout := widgets.NewQGridLayout2()
	masterLayout.AddLayout(headerLayout, 0, 0, 0)
	//masterLayout.AddLayout(portLayout, 1, 0, 0)
	masterLayout.AddLayout(infoLayout, 1, 0, 0)
	masterLayout.AddLayout(descriptionLayout, 2, 0, 0)
	masterLayout.AddLayout(bottomButtons, 3, 0, 0)

	centralWidget := widgets.NewQWidget(d.window, 0)
	centralWidget.SetLayout(masterLayout)
	d.window.SetCentralWidget(centralWidget)

	rxsmIcon := gui.NewQIcon5(GlobalConfig.IconPath)
	d.app.SetWindowIcon(rxsmIcon)

	d.window.Show()
	if len(d.saveHandler.PlaySaves) == 0 { // automatically prompt for RCX location if not default
		log.Println("PlaySaves is empty, opening install location dialog")
		d.installPathDialog = NewInstallPathDialog(d.window, 0)
		d.installPathDialog.ConnectFinished(d.onInstallPathDialogFinished)
		d.installPathDialog.OpenInstallPathDialog()
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-settings", "--settings":
			log.Println("Found -settings run arg, opening settingsDialog")
			d.settingsButton.Click()
		case "-install-path", "--install-path":
			if d.installPathDialog == nil {
				log.Println("Found -install-path run arg, opening installPathDialog")
				d.installPathDialog = NewInstallPathDialog(d.window, 0)
				d.installPathDialog.ConnectFinished(d.onInstallPathDialogFinished)
				d.installPathDialog.OpenInstallPathDialog()
			}
		case "-versioning", "--versioning":
			log.Println("Found -versions run arg, opening versionDialog")
			d.versionsButton.Click()
		}
	}
	go d.checkForUpdates()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	d.app.Exec()
	d.window.Hide()
	if d.saveVersioner != nil {
		d.saveVersioner.Exit()
	}
	log.Println("Display ended")
	d.endChan <- d.exitStatus
}

func (d *Display) Start() {
	go d.Run()
}

func (d *Display) Join() (int, error) {
	return <-d.endChan, nil
}

func (d *Display) populateFields() {
	if d.selectedSave == nil {
		return
	}
	d.nameField.SetText(d.selectedSave.Data.Name)
	d.creatorField.SetText(d.selectedSave.Data.Creator)
	d.idLabel.SetText("ID: " + DoubleDigitStr(d.selectedSave.Data.Id))
	d.descriptionField.SetPlainText(d.selectedSave.Data.Description)
	d.thumbnailImage.Swap(gui.NewQIcon5(d.selectedSave.ThumbnailPath()))
	d.thumbnailButton.SetIcon(d.thumbnailImage)
	if d.activeSave != nil && d.selectedSave != nil && *d.activeSave == *d.selectedSave {
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
	if d.temporaryThumbnailPath != d.selectedSave.ThumbnailPath() && d.temporaryThumbnailPath != "" {
		copyErr := d.saveHandler.CopyTo(d.temporaryThumbnailPath, d.selectedSave.ThumbnailPath())
		if copyErr != nil {
			log.Print("Error while copying Thumbnail from " + d.temporaryThumbnailPath + " to " + d.selectedSave.ThumbnailPath())
			log.Println(copyErr)
		}
		d.temporaryThumbnailPath = ""
	}
}

// start event methods

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
	d.onResetFilterButtonClicked(false)
	log.Println("Switched to mode " + strconv.Itoa(d.activeMode))
}

func (d *Display) onSaveSelectedChanged(index int) {
	if index == -1 { // no items in dropdown
		return
	}
	d.selectedSave = &(*d.activeSaves)[index]
	d.populateFields()
	log.Println("Selected " + strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onCopySaveButtonClicked(bool) {
	// very similar to onNewSaveButtonClicked
	newId := d.saveHandler.MaxId() + 1
	var dupSave Save
	var dupSaveErr error
	switch d.activeMode {
	case BUILD_MODE:
		newFolder := d.saveHandler.BuildSaveFolderPath(newId)
		dupSave, dupSaveErr = d.selectedSave.Duplicate(newFolder, newId)
		if dupSaveErr != nil {
			log.Println("Error while duplicating save")
			log.Println(dupSaveErr)
			return
		}
		d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves, dupSave)
	case PLAY_MODE:
		newFolder := d.saveHandler.PlaySaveFolderPath(newId)
		dupSave, dupSaveErr = d.selectedSave.Duplicate(newFolder, newId)
		if dupSaveErr != nil {
			log.Println("Error while duplicating save")
			log.Println(dupSaveErr)
			return
		}
		d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves, dupSave)
	}
	d.saveSelector.AddItems(makeSelectorOptions([]Save{dupSave}))
	log.Println("Copied save " + strconv.Itoa(d.selectedSave.Data.Id) + " to " + strconv.Itoa(dupSave.Data.Id))
	// select copied save
	d.saveSelector.SetCurrentIndex(len(*d.activeSaves) - 1)
	// propagation calls d.onSaveSelectedChanged(len(*d.activeSaves)-1)
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
	d.saveSelector.AddItems(makeSelectorOptions([]Save{newSave}))
	log.Println("Created new save " + strconv.Itoa(newSave.Data.Id))
	// select newly created save
	d.saveSelector.SetCurrentIndex(len(*d.activeSaves) - 1)
	// propagation calls d.onSaveSelectedChanged(len(*d.activeSaves)-1)
}

func (d *Display) onApplyFilterButtonClicked(bool) {
	log.Println("Applying save filter")
	isMatch := func (search string, s Save) bool{return false}
	switch d.filterTypeSelector.CurrentIndex(){
	case 0: // "Any"
		isMatch = isAnyMatch
	case 1: // "Name"
		isMatch = isNameMatch
	case 2: // "Creator"
		isMatch = isCreatorMatch
	case 3: // "Description"
		isMatch = isDescriptionMatch
	case 4: // "ID"
		isMatch = isIDMatch
	}
	resultChan := make(chan int)
	doMatch := func (search string, s Save, arrayLoc int) {
		if isMatch(search, s) {
			resultChan <- arrayLoc
		} else {
			resultChan <- -1
		}
	}
	searchString := d.filterSearchField.Text()
	for i, save := range d.filteredSaves {
		go doMatch(searchString, save, i)
	}
	oldFilterLen := len(d.filteredSaves)
	newFilteredSaves := []Save{}
	newFilterMapping := make(map[int]int)
	for j:=0; j < oldFilterLen; j++ {
		oldLoc := <- resultChan
		if oldLoc != -1 {
			newFilteredSaves = append(newFilteredSaves, d.filteredSaves[oldLoc])
			newFilterMapping[len(newFilteredSaves)-1] = d.filterMapping[oldLoc]
		}
	}
	d.filteredSaves = newFilteredSaves
	d.filterMapping = newFilterMapping
	// refill dropdown
	d.filteredSaveSelector.Clear()
	d.filteredSaveSelector.AddItems(makeSelectorOptions(d.filteredSaves))
}

func (d *Display) onResetFilterButtonClicked(bool) {
	log.Println("Reseting save filter")
	d.filteredSaves = *d.activeSaves
	// build mapping
	d.filterMapping = make(map[int]int)
	for i:=0; i < len(d.filteredSaves); i++ {
		d.filterMapping[i] = i
	}
	// refill dropdown
	d.filteredSaveSelector.Clear()
	d.filteredSaveSelector.AddItems(makeSelectorOptions(d.filteredSaves))
}

func (d *Display) onFilterSaveSelectedChanged(index int) {
	d.saveSelector.SetCurrentIndex(d.filterMapping[index])
}

func (d *Display) onSettingsButtonClicked(bool) {
	log.Println("Opening settings window")
	if d.settingsDialog == nil {
		d.settingsDialog = NewSettingsDialog(d.window, 0)
		d.settingsDialog.ConnectFinished(d.onSettingsDialogFinished)
	}
	d.settingsDialog.OpenSettingsDialog()
}

func (d *Display) onSettingsDialogFinished(i int) {
	// Nothing happens here
	log.Println("Settings window closed with code " + strconv.Itoa(i))
}

func (d *Display) onImportButtonClicked(bool) {
	var fileDialog *widgets.QFileDialog = widgets.NewQFileDialog(nil, 0)
	importPath := fileDialog.GetOpenFileName(d.window, "Select the file to import", "", "Zip Archive (*.zip);;Any File (*.*)", "", 0)
	var importedSaves []*Save
	var importErr error
	if importPath != "" {
		importStart := time.Now()
		switch d.activeMode {
		case BUILD_MODE:
			importedSaves, importErr = Import(importPath, d.saveHandler.BuildPath())
			if importErr != nil {
				log.Println("Import from " + importPath + " to " + d.saveHandler.BuildPath() + " failed")
				log.Println(importErr)
				return
			}
			for _, save := range importedSaves {
				d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves, *save)
				d.saveSelector.AddItems(makeSelectorOptions([]Save{*save}))
			}
		case PLAY_MODE:
			importedSaves, importErr = Import(importPath, d.saveHandler.PlayPath())
			if importErr != nil {
				log.Println("Import from " + importPath + " to " + d.saveHandler.PlayPath() + " failed")
				log.Println(importErr)
				return
			}
			for _, save := range importedSaves {
				d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves, *save)
				d.saveSelector.AddItems(makeSelectorOptions([]Save{*save}))
			}
		}
		log.Println("Imported " + strconv.Itoa(len(importedSaves)) + " saves in " + strconv.FormatFloat(time.Since(importStart).Seconds(), 'f', -1, 64) + "s")
	}
}

func (d *Display) onExportButtonClicked(bool) {
	var fileDialog *widgets.QFileDialog = widgets.NewQFileDialog(nil, 0)
	exportPath := fileDialog.GetSaveFileName(d.window, "Select the export location", "", "Zip Archive (*.zip)", "", 0)
	if exportPath != "" {
		exportErr := Export(exportPath, *d.selectedSave)
		if exportErr != nil {
			log.Println("Export to " + exportPath + " failed for " + strconv.Itoa(d.selectedSave.Data.Id))
			log.Println(exportErr)
			return
		} else {
			log.Println("Exported save " + strconv.Itoa(d.selectedSave.Data.Id) + " to " + exportPath)
		}
	} else {
		log.Println("Export file dialog dismissed")
	}
}

func (d *Display) onThumbnailButtonClicked(bool) {
	var fileDialog *widgets.QFileDialog = widgets.NewQFileDialog(nil, 0)
	d.temporaryThumbnailPath = fileDialog.GetOpenFileName(d.window, "Select a new thumbnail", d.selectedSave.ThumbnailPath(), "Images (*.jpg)", "", 0)
	if d.temporaryThumbnailPath != "" {
		d.thumbnailImage.Swap(gui.NewQIcon5(d.temporaryThumbnailPath))
		d.thumbnailButton.SetIcon(d.thumbnailImage)
		log.Println("Thumbnail temporarily set to " + d.temporaryThumbnailPath)
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
	log.Println("Saved " + strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onCancelButtonClicked(bool) {
	d.populateFields()
	log.Println("Canceled " + strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onActivateChecked(checkState int) {
	if d.activeMode == PLAY_MODE {
		log.Println("Hey buddy, you can't activate a play save")
		return // activating a play save is pointless
	}
	switch checkState {
	case 0:
		if d.activeSave == nil || *d.activeSave != *d.selectedSave {
			// automatic de-check when save changed or no currently active save
			return
		}
		moveErr := d.activeSave.MoveToId()
		if moveErr != nil {
			log.Println("Error while deactivating " + strconv.Itoa(d.activeSave.Data.Id))
			log.Println(moveErr)
			return
		}
		if d.saveVersioner != nil {
			d.saveVersioner.Exit()
		}
		d.saveVersioner = nil
		log.Println("Deactivated " + strconv.Itoa(d.activeSave.Data.Id))
		d.activeSave = nil
	case 2:
		if d.activeSave != nil && d.selectedSave != nil && *d.activeSave == *d.selectedSave && d.saveVersioner != nil {
			// automatic re-check when active save selected
			return
		}
		if d.activeSave != nil {
			// deactivate old activate save
			moveErr := d.activeSave.MoveToId()
			if moveErr != nil {
				log.Println("Error while deactivating " + strconv.Itoa(d.activeSave.Data.Id))
				log.Println(moveErr)
				return
			}
		}
		if d.selectedSave != nil {
			// activate new active save
			d.activeSave = d.selectedSave
			moveErr := d.activeSave.MoveToFirst()
			if moveErr != nil {
				log.Println("Error while activating " + strconv.Itoa(d.activeSave.Data.Id))
				log.Println(moveErr)
				return
			}
			if d.saveVersioner != nil {
				d.saveVersioner.Exit()
			}
			sv, svErr := NewSaveVersioner(d.activeSave)
			if svErr != nil {
				log.Println("Error creating SaveVersioner for save " + strconv.Itoa(d.activeSave.Data.Id))
				log.Println(svErr)
				return
			}
			d.saveVersioner = sv
			d.saveVersioner.Start(GlobalConfig.SnapshotPeriod)
			log.Println("Activated " + strconv.Itoa(d.activeSave.Data.Id))
		} else {
			log.Println("Selected save is nil, activation failed")
		}
	}
}

func (d *Display) onMoveToButtonClicked(bool) {
	if d.selectedSave == nil {
		return
	}
	log.Println("Moving save " + strconv.Itoa(d.selectedSave.Data.Id))
	if d.selectedSave == d.activeSave {
		d.activeSave = nil
		d.selectedSave.MoveToId()
	}
	switch d.activeMode {
	case BUILD_MODE:
		moveErr := d.selectedSave.Move(d.saveHandler.PlaySaveFolderPath(d.selectedSave.Data.Id))
		if moveErr != nil {
			log.Println("Error while moving " + strconv.Itoa(d.selectedSave.Data.Id))
			log.Println(moveErr)
			return
		}
		selIndex := d.saveSelector.CurrentIndex()
		d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves, d.saveHandler.BuildSaves[selIndex])                    // add to playsaves
		d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves[:selIndex], d.saveHandler.BuildSaves[selIndex+1:]...) // remove from buildsaves
		d.modeTab.SetCurrentIndex(PLAY_MODE)                                                                             // toggle to other mode to keep showing selected save
		// propagation calls d.onModeTabChanged(PLAY_MODE)
	case PLAY_MODE:
		moveErr := d.selectedSave.Move(d.saveHandler.BuildSaveFolderPath(d.selectedSave.Data.Id))
		if moveErr != nil {
			log.Println("Error while moving " + strconv.Itoa(d.selectedSave.Data.Id))
			log.Println(moveErr)
			return
		}
		selIndex := d.saveSelector.CurrentIndex()
		d.saveHandler.BuildSaves = append(d.saveHandler.BuildSaves, d.saveHandler.PlaySaves[selIndex])                // add to playsaves
		d.saveHandler.PlaySaves = append(d.saveHandler.PlaySaves[:selIndex], d.saveHandler.PlaySaves[selIndex+1:]...) // remove from buildsaves
		d.modeTab.SetCurrentIndex(BUILD_MODE)                                                                         // toggle to other mode to keep showing selected save
		// propagation calls d.onModeTabChanged(BUILD_MODE)
	}
	if d.activeSave == nil && len(*d.activeSaves) > 0 {
		d.activeSave = &(*d.activeSaves)[0]
		d.activeSave.MoveToFirst()
	}
	// re-select save
	d.saveSelector.SetCurrentIndex(len(*d.activeSaves) - 1)
	//d.onSaveSelectedChanged(len(*d.activeSaves)-1)
	log.Println("Save moved to " + strconv.Itoa(d.activeMode))
}

func (d *Display) onVersionsButtonClicked(bool) {
	if d.saveVersioner == nil {
		log.Println("save versioner is nil, ignoring version button click")
		return
	}
	d.saveVersioner.Exit()
	versionDialog := NewVersionDialog(d.window, 0)
	versionDialog.ConnectFinished(d.onVersionsDialogFinished)
	versionDialog.OpenVersionDialog(d.saveVersioner)
}

func (d *Display) onVersionsDialogFinished(int) {
	GlobalConfig.Save()
	d.saveVersioner.Start(GlobalConfig.SnapshotPeriod)
	// TODO: reload selectedSave
	d.activeSave.FreeID()
	newS, err := NewSave(d.activeSave.FolderPath())
	if err != nil {
		log.Println("Error reloading selected save file")
		log.Println(err)
	}
	*d.activeSave = newS
	if d.selectedSave != nil && *d.activeSave == *d.selectedSave {
		d.onSaveSelectedChanged(d.saveSelector.CurrentIndex())
	}
}

// end event methods

func (d *Display) checkForUpdates() {
	if !GlobalConfig.AutoCheck {
		return
	}
	log.Println("Checking for updates")
	_, _, ok := checkForRXSMUpdate()
	if !ok {
		log.Println("Checking for updates failed")
		return
	}
	log.Println("Update available", IsOutOfDate)
	if IsOutOfDate {
		log.Println("New update download link " + DownloadURL)
		if GlobalConfig.AutoInstall {
			downloadRXSMUpdateQuiet()
		}
	}
}

// end Display

func makeSelectorOptions(saves []Save) []string {
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
