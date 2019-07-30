// Created 2019-07-26 by NGnius

package display

import (
	"strconv"
	"log"
	"os"

	"github.com/therecipe/qt/widgets"

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
	selectedSave *saver.Save
	activeSave *saver.Save
	saveHandler saver.SaveHandler
	endChan chan int
	// Qt GUI objects
	window *widgets.QMainWindow
	app *widgets.QApplication
	/* TODO: import and export buttons + functionality
	importButton *widget.QPushButton2
	exportButton *widget.QPushButton2
	*/
	saveSelector *widgets.QComboBox
	nameField *widgets.QLineEdit
	creatorLabel *widgets.QLabel
	creatorField *widgets.QLineEdit
	/* TODO: implement thumbnail button + functionality
	imageLabel *widgets.QPixmap
	imageButton *widgets.QPushButton
	*/
	idLabel *widgets.QLabel
	descriptionLabel *widgets.QLabel
	descriptionField *widgets.QPlainTextEdit
	saveButton *widgets.QPushButton
	cancelButton *widgets.QPushButton
	activateButton *widgets.QPushButton
	moveButton *widgets.QPushButton
}

func NewDisplay(saveHandler saver.SaveHandler) (*Display){
	newD := Display {endChan: make(chan int), saveHandler:saveHandler}
	// set to invalid Id
	newD.activeSave = saveHandler.ActiveBuildSave()
	return &newD
}

func (d *Display) Run() {
	log.Println("Display started")
	// build initial display
	d.app = widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	// with a minimum size of 250*200
	// and sets the title to "Hello Widgets Example"
	d.window = widgets.NewQMainWindow(nil, 0)
	d.window.SetMinimumSize2(250, 200)
	d.window.SetWindowTitle("rxsm")

	d.saveSelector = widgets.NewQComboBox(nil)
	d.saveSelector.AddItems(makeSelectorOptions(d.saveHandler.BuildSaves))
	d.saveSelector.ConnectCurrentIndexChanged(d.onSaveSelectedChanged)

	d.nameField = widgets.NewQLineEdit(nil)
	d.creatorLabel = widgets.NewQLabel2("by", nil, 0)
	d.creatorField = widgets.NewQLineEdit(nil)
	d.idLabel = widgets.NewQLabel2("id: ##", nil, 0)
	d.descriptionLabel = widgets.NewQLabel2("Description:", nil, 0)
	d.descriptionField = widgets.NewQPlainTextEdit(nil)

	d.saveButton = widgets.NewQPushButton2("save", nil)
	d.saveButton.ConnectClicked(d.onSaveButtonClicked)
	d.cancelButton = widgets.NewQPushButton2("cancel", nil)
	d.cancelButton.ConnectClicked(d.onCancelButtonClicked)
	d.activateButton = widgets.NewQPushButton2("activate", nil)
	d.activateButton.ConnectClicked(d.onActivateButtonClicked)
	d.moveButton = widgets.NewQPushButton2("functional button", nil)
	d.moveButton.ConnectClicked(d.onMoveToButtonClicked)

	d.selectedSave = &d.saveHandler.BuildSaves[d.saveSelector.CurrentIndex()]
	d.populateFields()

	headerLayout := widgets.NewQGridLayout2()
	headerLayout.AddWidget2(d.saveSelector, 0, 0, 0)

	infoLayout := widgets.NewQGridLayout2()
	infoLayout.AddWidget3(d.nameField, 0, 0, 1, 6, 0)
	infoLayout.AddWidget2(d.creatorLabel, 1, 0, 0)
	infoLayout.AddWidget3(d.creatorField, 1, 1, 1, 5, 0)

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
	d.nameField.SetText(d.selectedSave.Data.Name)
	d.creatorField.SetText(d.selectedSave.Data.Creator)
	oldIdText := d.idLabel.Text()
	d.idLabel.SetText(oldIdText[:len(oldIdText)-2]+saver.DoubleDigitStr(d.selectedSave.Data.Id))
	d.descriptionField.SetPlainText(d.selectedSave.Data.Description)
}

func (d *Display) syncBackFields() {
	d.selectedSave.Data.Name = d.nameField.Text()
	d.selectedSave.Data.Creator = d.creatorField.Text()
	d.selectedSave.Data.Description = d.descriptionField.ToPlainText()
}

func (d *Display) onSaveSelectedChanged(index int) {
	d.selectedSave = &d.saveHandler.BuildSaves[index]
	d.populateFields()
	log.Println("Selected "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onSaveButtonClicked(bool) {
	d.syncBackFields()
	saveErr := d.selectedSave.Data.Save()
	if saveErr != nil {
		log.Println(saveErr)
	}
	index := d.saveSelector.CurrentIndex()
	d.saveSelector.SetItemText(index, makeSelectorOptions(d.saveHandler.BuildSaves)[index])
	log.Println("Saved "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onCancelButtonClicked(bool) {
	d.populateFields()
	log.Println("Canceled "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onActivateButtonClicked(bool) {
	moveSaveToFirst(d.selectedSave, d.saveHandler.BuildSaves)
	log.Println("Activated "+strconv.Itoa(d.selectedSave.Data.Id))
}

func (d *Display) onMoveToButtonClicked(bool) {
	// TODO: implement move to opposite build/play game mode folder
	log.Println("Move to button clicked (unimplemented) button clicked")
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

func moveSaveToFirst(selected *saver.Save, saves []saver.Save) {
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
