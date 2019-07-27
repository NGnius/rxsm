// Created 2019-07-26 by NGnius

package display

import (
	"fyne.io/fyne/widget"
	"fyne.io/fyne/app"

  "../saver"
)

func Run(saveHandler saver.SaveHandler) {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()
}
