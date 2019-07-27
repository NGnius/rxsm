// Created 2019-07-26 by NGnius

package display

import (
	//"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/app"
	//"strconv"
	"log"
	//"strings"

  "../saver"
)

var (
	Width int = 640
	Height int = 480
	savesSelector widget.Select
	activeSaveHandler saver.SaveHandler
	saveGroup widget.Group
	gameIdEntry *widget.Entry = widget.NewEntry()
	gameNameEntry *widget.Entry = widget.NewEntry()
	gameDescriptionEntry *widget.Entry = widget.NewMultiLineEntry()
	gameCreatorEntry *widget.Entry = widget.NewEntry()
	gameId widget.FormItem = widget.FormItem{Text:"Game Name", Widget:gameIdEntry}
	gameName widget.FormItem = widget.FormItem{Text:"Game Name", Widget:gameNameEntry}
	gameDescription widget.FormItem = widget.FormItem{Text:"Game Name", Widget:gameDescriptionEntry}
	gameCreator widget.FormItem = widget.FormItem{Text:"Game Name", Widget:gameCreatorEntry}
	saveChangesButton *widget.Button = widget.NewButton("Save Changes", onSaveButtonClick)
)

func Run(saveHandler saver.SaveHandler) {
	activeSaveHandler = saveHandler
	app := app.New()

	savesSelector := widget.NewSelect(
		makeSelectorOptions(saveHandler.BuildSaves),
		onSelectionChange)

	dataForm := widget.NewForm(&gameId, &gameName, &gameDescription)

	saveGroup := widget.NewGroup("Save", savesSelector, dataForm, saveChangesButton)

	w := app.NewWindow("rxsm")
	w.SetContent(widget.NewVBox(
		saveGroup,
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	//w.Resize(fyne.NewSize(Width, Height))
	w.ShowAndRun()
}

// start UI interaction events
func onSelectionChange(selectedOption string) {
	log.Println("Selector changed to: "+selectedOption)
}

func onSaveButtonClick() {
	// save changes to GameData
	save, ok := getSelectedSave(savesSelector.Selected)
	if !ok { // nothing (or invalid item) is selected
		return
	}
	//save.Data.Id = gameId.Text
	save.Data.Name = gameNameEntry.Text
	save.Data.Description = gameDescriptionEntry.Text
	save.Data.Creator = gameCreatorEntry.Text
	saveErr := save.Data.Save()
	if saveErr != nil {
		log.Println("Error while saving: ")
		log.Println(saveErr)
	}
}
// end UI events


func makeSelectorOptions(saves []saver.Save) ([]string) {
	var result []string
	for _, s := range saves {
		result = append(result, s.Data.Name)
	}
	log.Println(result)
	return result
}

func moveSaveToFirst(selected string, saves []saver.Save) (ok bool) {
	for _, s := range saves {
		err := s.MoveOut()
		if err != nil {
			log.Println(err)
		}
	}
	save, ok := getSelectedSave(selected)
	if !ok {
		log.Println("Failed to find selected save!")
		return false
	}
	save.MoveToFirst()
	return true
}

func getSelectedSave(name string) (save saver.Save, isFound bool){
	var noResult saver.Save
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
