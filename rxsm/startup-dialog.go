// Created 2019-08-18 by NGnius

package main

import (
  //"log"
  //"os"

  "github.com/therecipe/qt/widgets"
)

// start InstallPathDialog

type InstallPathDialog struct {
  widgets.QDialog
  InstallPath string
  installPathChan chan string
  Cancelled bool
  infoLabel *widgets.QLabel
  pathField *widgets.QLineEdit
  browseButton *widgets.QPushButton
  fillerLabel *widgets.QLabel
  okButton *widgets.QPushButton
  cancelButton *widgets.QPushButton
}

// NewInstallPathDialog(parent *widgets.QWidget, flags) is automatically generated

func (ipd *InstallPathDialog) OpenInstallPathDialog() (installPath string) {
  // TODO
  ipd.Cancelled = true // assume cancelled unless proven otherwise
  ipd.__init_display()
  ipd.Open()
  return ipd.InstallPath
}

func (ipd *InstallPathDialog) __init_display() {
  ipd.installPathChan = make(chan string, 1)
  // build dialog window
  ipd.infoLabel = widgets.NewQLabel2("Unable to find your RobocraftX saves. <b>Please specify where RCX is installed.</b> <br/>For advanced configuration, see the <a href='https://github.com/NGnius/rxsm/wiki/User-Guide#configuration'>User Guide<a>", nil, 0)
  ipd.infoLabel.SetTextFormat(1) // rich text (html subset)
  ipd.pathField = widgets.NewQLineEdit(nil)
  ipd.browseButton = widgets.NewQPushButton2("Browse", nil)
  ipd.browseButton.ConnectClicked(ipd.onBrowseButtonClicked)

  ipd.fillerLabel = widgets.NewQLabel2("To apply this change, click Ok and restart RXSM", nil, 0)
  ipd.okButton = widgets.NewQPushButton2("Ok", nil)
  ipd.okButton.ConnectClicked(ipd.onOkButtonClicked)
  ipd.cancelButton = widgets.NewQPushButton2("Cancel", nil)
  ipd.cancelButton.ConnectClicked(ipd.onCancelButtonClicked)

  infoLayout := widgets.NewQGridLayout2()
  infoLayout.AddWidget2(ipd.infoLabel, 0, 0, 0)
	infoLayout.AddWidget3(ipd.pathField, 1, 0, 1, 4, 0)
  infoLayout.AddWidget2(ipd.browseButton, 1, 4, 0)

  confirmLayout := widgets.NewQGridLayout2()
  confirmLayout.AddWidget3(ipd.fillerLabel, 0, 0, 1, 3, 0)
  confirmLayout.AddWidget2(ipd.okButton, 0, 3, 0)
  confirmLayout.AddWidget2(ipd.cancelButton, 0, 4, 0)

  masterLayout := widgets.NewQGridLayout2()
  masterLayout.AddLayout(infoLayout, 0, 0, 0)
  masterLayout.AddLayout(confirmLayout, 1, 0, 0)

  ipd.SetLayout(masterLayout)
}

func (ipd *InstallPathDialog) onBrowseButtonClicked(bool) {
  // TODO: open folder picker dialog
  var fileDialog *widgets.QFileDialog = widgets.NewQFileDialog(nil, 0)
	ipd.pathField.SetText(fileDialog.GetExistingDirectory(ipd, "Select a folder", "", 0))
}

func (ipd *InstallPathDialog) onOkButtonClicked(bool) {
  ipd.Cancelled = false
  ipd.InstallPath = ipd.pathField.Text()
  ipd.Accept()
}

func (ipd *InstallPathDialog) onCancelButtonClicked(bool) {
  ipd.Cancelled = true // assumed already
  ipd.InstallPath = ""
  ipd.Reject()
}

// end InstallPathDialog
