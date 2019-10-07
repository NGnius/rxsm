// Created 2019-09-04 by NGnius

package main

import (
	//"log"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// start SettingsDialog

type SettingsDialog struct {
	widgets.QDialog
	isDisplayInited bool
	// top
	settingsLabel *widgets.QLabel
	tabs          *widgets.QTabWidget
	saveSettings  *widgets.QWidget
	rxsmSettings  *widgets.QWidget
	aboutSettings *widgets.QWidget // arguably not settings
	// save settings widgets
	saveLabel           *widgets.QLabel
	creatorLabel        *widgets.QLabel
	creatorField        *widgets.QLineEdit
	forceCreatorLabel   *widgets.QLabel
	forceCreatorField   *widgets.QCheckBox
	forceUniqueIdsLabel *widgets.QLabel
	forceUniqueIdsField *widgets.QCheckBox
	defaultSaveLabel    *widgets.QLabel
	defaultSaveField    *widgets.QLineEdit
	advancedLabel       *widgets.QLabel
	playLabel           *widgets.QLabel
	playField           *widgets.QLineEdit
	buildLabel          *widgets.QLabel
	buildField          *widgets.QLineEdit
	// rxsm config widgets
	configLabel         *widgets.QLabel
	logLabel            *widgets.QLabel
	logField            *widgets.QLineEdit
	appIconLabel        *widgets.QLabel
	appIconField        *widgets.QLineEdit
	iconPackLabel       *widgets.QLabel
	iconPackField       *widgets.QLineEdit
	snapshotPeriodLabel *widgets.QLabel
	snapshotPeriodField *widgets.QLineEdit
	updateLabel         *widgets.QLabel
	autoCheckUpdate     *widgets.QCheckBox
	autoInstallUpdate   *widgets.QCheckBox
	doNotTrack          *widgets.QCheckBox
	updateServerLabel   *widgets.QLabel
	updateServerField   *widgets.QLineEdit
	rxsmFiller          *widgets.QLabel
	// about widgets
	iconLabel         *widgets.QLabel
	rxsmVersionLabel  *widgets.QLabel
	machineLabel      *widgets.QLabel
	descriptionLabel  *widgets.QLabel
	updateButton      *widgets.QPushButton
	updateProgressBar *widgets.QProgressBar
	// bottom
	fillerLabel  *widgets.QLabel
	okButton     *widgets.QPushButton
	cancelButton *widgets.QPushButton
}

// NewSettingsDialog(parent *widgets.QWidget, flags) is automatically generated

func (sd *SettingsDialog) OpenSettingsDialog() {
	if !sd.isDisplayInited {
		sd.__init_display()
	}
	sd.populateFields()
	sd.Open()
}

func (sd *SettingsDialog) __init_display() {
	// top
	sd.settingsLabel = widgets.NewQLabel2("<b>RXSM Settings & Configuration</b> <br/>For more info, see <a href='https://github.com/NGnius/rxsm/wiki/User-Guide#advanced-configuration'>Advanced Configuration</a> <br/><i>Some values require a restart to take effect</i>", nil, 0)
	sd.tabs = widgets.NewQTabWidget(nil)
	sd.saveSettings = widgets.NewQWidget(nil, 0)
	sd.rxsmSettings = widgets.NewQWidget(nil, 0)
	sd.aboutSettings = widgets.NewQWidget(nil, 0)
	sd.tabs.AddTab(sd.saveSettings, "Save Settings")
	sd.tabs.AddTab(sd.rxsmSettings, "Configuration")
	sd.tabs.AddTab(sd.aboutSettings, "About")

	topLayout := widgets.NewQGridLayout2()
	topLayout.AddWidget2(sd.settingsLabel, 0, 0, 0)
	topLayout.AddWidget2(sd.tabs, 1, 0, 0)
	topLayout.SetRowStretch(1, 1)

	masterLayout := widgets.NewQGridLayout2()
	masterLayout.AddLayout(topLayout, 0, 0, 0)

	// save settings tab
	sd.saveLabel = widgets.NewQLabel2("<b>Settings for RXSM management of RobocraftX saves</b>", nil, 0)
	sd.creatorLabel = widgets.NewQLabel2("Creator", nil, 0)
	sd.creatorField = widgets.NewQLineEdit(nil)
	sd.creatorField.SetToolTip("The value to use as 'CreatorName' for new game saves")
	sd.creatorLabel.SetBuddy(sd.creatorField)
	sd.forceCreatorLabel = widgets.NewQLabel2("Force Creator", nil, 0)
	sd.forceCreatorField = widgets.NewQCheckBox2("", nil)
	sd.forceCreatorField.SetToolTip("Check this to force 'Creator' for old game saves as well")
	sd.forceCreatorLabel.SetBuddy(sd.forceCreatorField)
	sd.forceUniqueIdsLabel = widgets.NewQLabel2("Force Unique IDs", nil, 0)
	sd.forceUniqueIdsField = widgets.NewQCheckBox2("", nil)
	sd.forceUniqueIdsField.SetToolTip("Check this to force game saves to have unique IDs")
	sd.forceUniqueIdsLabel.SetBuddy(sd.forceUniqueIdsField)
	sd.defaultSaveLabel = widgets.NewQLabel2("Default Save", nil, 0)
	sd.defaultSaveField = widgets.NewQLineEdit(nil)
	sd.defaultSaveField.SetToolTip("The folder of the game save to copy when 'New' is clicked")
	sd.defaultSaveLabel.SetBuddy(sd.defaultSaveField)

	sd.advancedLabel = widgets.NewQLabel2("&nbsp;&nbsp;&nbsp;&nbsp;<b>Advanced</b>", nil, 0)
	sd.playLabel = widgets.NewQLabel2("Play Path", nil, 0)
	sd.playField = widgets.NewQLineEdit(nil)
	sd.playField.SetToolTip("The folder directly containing all community game save folders")
	sd.playLabel.SetBuddy(sd.playField)
	sd.buildLabel = widgets.NewQLabel2("Build Path", nil, 0)
	sd.buildField = widgets.NewQLineEdit(nil)
	sd.buildField.SetToolTip("The folder directly containing all creative game save folders")
	sd.buildLabel.SetBuddy(sd.buildField)

	saveLayout := widgets.NewQGridLayout2()
	saveLayout.AddWidget3(sd.saveLabel, 0, 0, 1, 3, 0)
	saveLayout.AddWidget2(sd.creatorLabel, 1, 0, 0)
	saveLayout.AddWidget3(sd.creatorField, 1, 1, 1, 2, 0)
	saveLayout.AddWidget2(sd.forceCreatorLabel, 2, 0, 0)
	saveLayout.AddWidget3(sd.forceCreatorField, 2, 1, 1, 2, 0)
	saveLayout.AddWidget2(sd.forceUniqueIdsLabel, 3, 0, 0)
	saveLayout.AddWidget3(sd.forceUniqueIdsField, 3, 1, 1, 2, 0)
	saveLayout.AddWidget2(sd.defaultSaveLabel, 4, 0, 0)
	saveLayout.AddWidget3(sd.defaultSaveField, 4, 1, 1, 2, 0)
	saveLayout.AddWidget3(sd.advancedLabel, 5, 0, 1, 3, 0)
	saveLayout.AddWidget2(sd.playLabel, 6, 0, 0)
	saveLayout.AddWidget3(sd.playField, 6, 1, 1, 2, 0)
	saveLayout.AddWidget2(sd.buildLabel, 7, 0, 0)
	saveLayout.AddWidget3(sd.buildField, 7, 1, 1, 2, 0)
	sd.saveSettings.SetLayout(saveLayout)

	// rxsm settings tab
	sd.configLabel = widgets.NewQLabel2("<b>Configurable Values for RXSM</b>", nil, 0)
	sd.logLabel = widgets.NewQLabel2("Log Path", nil, 0)
	sd.logField = widgets.NewQLineEdit(nil)
	sd.logField.SetToolTip("The file to write log events to (all parent folders must exist already)")
	sd.logLabel.SetBuddy(sd.logField)
	sd.appIconLabel = widgets.NewQLabel2("App Logo Path", nil, 0)
	sd.appIconField = widgets.NewQLineEdit(nil)
	sd.appIconField.SetToolTip("The icon file (.svg or .jpg) to use as RXSM's logo")
	sd.appIconLabel.SetBuddy(sd.appIconField)
	sd.iconPackLabel = widgets.NewQLabel2("Icon Folder Path", nil, 0)
	sd.iconPackField = widgets.NewQLineEdit(nil)
	sd.iconPackField.SetToolTip("The icon folder containing icon files (.svg or .jpg) for RXSM to use ")
	sd.iconPackLabel.SetBuddy(sd.iconPackField)
	sd.snapshotPeriodLabel = widgets.NewQLabel2("Snapshot Period (ns)", nil, 0)
	sd.snapshotPeriodField = widgets.NewQLineEdit(nil)
	sd.snapshotPeriodField.SetToolTip("The time (in nanoseconds) between automatic snapshots of the active save (0=disable)")
	sd.snapshotPeriodLabel.SetBuddy(sd.snapshotPeriodField)
	intValidator := gui.NewQIntValidator(nil)
	intValidator.SetBottom(0)
	sd.snapshotPeriodField.SetValidator(intValidator)
	sd.updateLabel = widgets.NewQLabel2("&nbsp;&nbsp;&nbsp;&nbsp;<b>Updates</b>", nil, 0)
	sd.autoCheckUpdate = widgets.NewQCheckBox2("Auto-Check", nil)
	sd.autoCheckUpdate.SetToolTip("Check this to automatically check for updates")
	sd.autoInstallUpdate = widgets.NewQCheckBox2("Auto-Update", nil)
	sd.autoInstallUpdate.SetToolTip("Check this to automatically install updates")
	sd.doNotTrack = widgets.NewQCheckBox2("D.N.T.", nil)
	sd.doNotTrack.SetToolTip("Check this to send a Do Not Track header with all web communication")
	sd.updateServerLabel = widgets.NewQLabel2("Update Server", nil, 0)
	sd.updateServerField = widgets.NewQLineEdit(nil)
	sd.updateServerField.SetToolTip("The server to use to check for updates")
	sd.updateServerLabel.SetBuddy(sd.updateServerField)
	sd.rxsmFiller = widgets.NewQLabel2("", nil, 0)

	configLayout := widgets.NewQGridLayout2()
	configLayout.AddWidget3(sd.configLabel, 0, 0, 1, 3, 0)
	configLayout.AddWidget2(sd.logLabel, 1, 0, 0)
	configLayout.AddWidget3(sd.logField, 1, 1, 1, 2, 0)
	configLayout.AddWidget2(sd.appIconLabel, 2, 0, 0)
	configLayout.AddWidget3(sd.appIconField, 2, 1, 1, 2, 0)
	configLayout.AddWidget2(sd.iconPackLabel, 3, 0, 0)
	configLayout.AddWidget3(sd.iconPackField, 3, 1, 1, 2, 0)
	configLayout.AddWidget2(sd.snapshotPeriodLabel, 4, 0, 0)
	configLayout.AddWidget3(sd.snapshotPeriodField, 4, 1, 1, 2, 0)
	configLayout.AddWidget3(sd.updateLabel, 5, 0, 1, 3, 0)
	configLayout.AddWidget2(sd.autoCheckUpdate, 6, 0, 0x0004)
	configLayout.AddWidget2(sd.autoInstallUpdate, 6, 1, 0x0004)
	configLayout.AddWidget2(sd.doNotTrack, 6, 2, 0x0004)
	configLayout.AddWidget2(sd.updateServerLabel, 7, 0, 0)
	configLayout.AddWidget3(sd.updateServerField, 7, 1, 1, 2, 0)
	// configLayout.AddWidget3(sd.rxsmFiller, 7, 0, 1, 3, 0)
	sd.rxsmSettings.SetLayout(configLayout)

	// about tab
	sd.iconLabel = widgets.NewQLabel2("", nil, 0)
	sd.iconLabel.SetAlignment(0x0084)
	logo := gui.NewQPixmap3(GlobalConfig.IconPath, "", 0).ScaledToHeight(80, 1)
	sd.iconLabel.SetPixmap(logo)
	sd.descriptionLabel = widgets.NewQLabel2("RobocraftX Save Manager, a <a href='https://github.com/NGnius/rxsm/blob/develop/LICENSE'>FOSS project</a> by NGnius to bring RCX players out of the Jurassic period. <br/><h3>RAWR!</h3>", nil, 0)
	sd.descriptionLabel.SetWordWrap(true)
	sd.descriptionLabel.SetAlignment(0x0004)
	versionStr := GlobalConfig.Version + " (" + runtime.Compiler + ")"
	if GlobalConfig.LastVersion() != GlobalConfig.Version && GlobalConfig.LastVersion() != "" {
		versionStr = GlobalConfig.LastVersion() + " -> " + versionStr
	}
	sd.rxsmVersionLabel = widgets.NewQLabel2("<b>Version</b> "+versionStr, nil, 0)
	sd.rxsmVersionLabel.SetAlignment(0x0004)
	sd.rxsmVersionLabel.SetSizePolicy2(1, 4)
	sd.machineLabel = widgets.NewQLabel2("<b>Machine</b> "+runtime.GOOS+"-"+runtime.GOARCH+" x"+strconv.Itoa(runtime.NumCPU())+" (go: "+strconv.Itoa(runtime.NumGoroutine())+")", nil, 0)
	sd.machineLabel.SetWordWrap(true)
	sd.machineLabel.SetAlignment(0x0004)
	sd.machineLabel.SetSizePolicy2(1, 4)
	sd.updateButton = widgets.NewQPushButton2("Check for Updates", nil)
	sd.updateButton.ConnectClicked(sd.onUpdateButtonClicked)
	sd.updateProgressBar = widgets.NewQProgressBar(nil)
	sd.updateProgressBar.SetTextVisible(true)
	sd.updateProgressBar.SetMaximum(UpdateSteps)

	aboutLayout := widgets.NewQGridLayout2()
	aboutLayout.AddWidget2(sd.iconLabel, 0, 0, 0)
	aboutLayout.AddWidget2(sd.descriptionLabel, 1, 0, 0)
	aboutLayout.AddWidget2(sd.rxsmVersionLabel, 2, 0, 0)
	aboutLayout.AddWidget2(sd.machineLabel, 3, 0, 0)
	aboutLayout.AddWidget2(sd.updateButton, 4, 0, 0)
	aboutLayout.AddWidget2(sd.updateProgressBar, 5, 0, 0)
	sd.updateProgressBar.Hide()
	sd.aboutSettings.SetLayout(aboutLayout)

	// bottom
	sd.fillerLabel = widgets.NewQLabel2("To apply changes, click Ok", nil, 0)
	sd.okButton = widgets.NewQPushButton2("Ok", nil)
	sd.okButton.ConnectClicked(sd.onOkButtonClicked)
	sd.cancelButton = widgets.NewQPushButton2("Cancel", nil)
	sd.cancelButton.ConnectClicked(sd.onCancelButtonClicked)

	bottomLayout := widgets.NewQGridLayout2()
	bottomLayout.AddWidget3(sd.fillerLabel, 0, 0, 1, 3, 0)
	bottomLayout.AddWidget2(sd.okButton, 0, 3, 0)
	bottomLayout.AddWidget2(sd.cancelButton, 0, 4, 0)
	masterLayout.AddLayout(bottomLayout, 1, 0, 0)
	sd.SetLayout(masterLayout)
	sd.isDisplayInited = true
}

func (sd *SettingsDialog) populateFields() {
	sd.populateSaveSettingsFields()
	sd.populateRXSMSettingsFields()
	sd.populateAboutFields()
}

func (sd *SettingsDialog) populateSaveSettingsFields() {
	sd.creatorField.SetText(GlobalConfig.Creator)
	creatorCheck := core.Qt__Unchecked
	if GlobalConfig.ForceCreator {
		creatorCheck = core.Qt__Checked
	}
	sd.forceCreatorField.SetCheckState(creatorCheck)
	idCheck := core.Qt__Unchecked
	if GlobalConfig.ForceUniqueIds {
		idCheck = core.Qt__Checked
	}
	sd.forceUniqueIdsField.SetCheckState(idCheck)
	sd.defaultSaveField.SetText(GlobalConfig.DefaultSaveFolder)
	sd.playField.SetText(GlobalConfig.PlayPath)
	sd.buildField.SetText(GlobalConfig.BuildPath)
}

func (sd *SettingsDialog) populateRXSMSettingsFields() {
	sd.logField.SetText(GlobalConfig.LogPath)
	sd.appIconField.SetText(GlobalConfig.IconPath)
	sd.iconPackField.SetText(GlobalConfig.IconPackPath)
	sd.snapshotPeriodField.SetText(strconv.Itoa(int(GlobalConfig.SnapshotPeriod)))
	autoCheckCheck := core.Qt__Unchecked
	if GlobalConfig.AutoCheck {
		autoCheckCheck = core.Qt__Checked
	}
	sd.autoCheckUpdate.SetCheckState(autoCheckCheck)
	autoInstallCheck := core.Qt__Unchecked
	if GlobalConfig.AutoInstall {
		autoInstallCheck = core.Qt__Checked
	}
	sd.autoInstallUpdate.SetCheckState(autoInstallCheck)
	dntCheck := core.Qt__Unchecked
	if GlobalConfig.DoNotTrack {
		dntCheck = core.Qt__Checked
	}
	sd.doNotTrack.SetCheckState(dntCheck)
	sd.updateServerField.SetText(GlobalConfig.UpdateServer)
}

func (sd *SettingsDialog) populateAboutFields() {
	if IsOutOfDate {
		sd.updateButton.SetEnabled(true)
		sd.updateButton.SetText("Install Update")
	} else if !IsUpdating {
		sd.updateButton.SetEnabled(true)
		sd.updateButton.SetText("Check for Updates")
	}
	if IsUpdating {
		sd.updateButton.SetEnabled(false)
		sd.updateButton.SetText("Updated (Restart Required)")
	}
}

func (sd *SettingsDialog) syncBackFields() {
	sd.syncBackSaveSettings()
	sd.syncBackRXSMSettings()
	GlobalConfig.Save()
}

func (sd *SettingsDialog) syncBackSaveSettings() {
	GlobalConfig.Creator = sd.creatorField.Text()
	GlobalConfig.ForceCreator = sd.forceCreatorField.IsChecked()
	GlobalConfig.ForceUniqueIds = sd.forceUniqueIdsField.IsChecked()
	GlobalConfig.DefaultSaveFolder = sd.defaultSaveField.Text()
	GlobalConfig.PlayPath = filepath.FromSlash(sd.playField.Text())
	GlobalConfig.BuildPath = filepath.FromSlash(sd.buildField.Text())
}

func (sd *SettingsDialog) syncBackRXSMSettings() {
	GlobalConfig.LogPath = filepath.FromSlash(sd.logField.Text())
	GlobalConfig.IconPath = filepath.FromSlash(sd.appIconField.Text())
	GlobalConfig.IconPackPath = filepath.FromSlash(sd.iconPackField.Text())
	newPeriod, parseErr := strconv.ParseInt(sd.snapshotPeriodField.Text(), 10, 64)
	if parseErr != nil {
		newPeriod = 0
	}
	GlobalConfig.SnapshotPeriod = newPeriod
	GlobalConfig.AutoCheck = sd.autoCheckUpdate.IsChecked()
	GlobalConfig.AutoInstall = sd.autoInstallUpdate.IsChecked()
	GlobalConfig.DoNotTrack = sd.doNotTrack.IsChecked()
	if GlobalConfig.DoNotTrack {
		ExtraHeader["DNT"] = []string{DNT_ON}
	} else {
		ExtraHeader["DNT"] = []string{DNT_OFF}
	}
	GlobalConfig.UpdateServer = sd.updateServerField.Text()
}

func (sd *SettingsDialog) onUpdateButtonClicked(bool) {
	// TODO: implement
	if !IsOutOfDate {
		go func() {
			sd.updateProgressBar.SetFormat("Checking for Update")
			sd.updateProgressBar.Show()
			_, _, ok := checkForRXSMUpdate()
			if !ok {
				return
			}
			if IsOutOfDate {
				sd.updateProgressBar.SetFormat("Update available")
				sd.updateButton.SetText("Install Update")
			}
			sd.updateProgressBar.Hide()
		}()
	} else if DownloadURL != "" {
		go func() {
			success := true
			sd.updateProgressBar.SetFormat("Updating - %p%")
			sd.updateProgressBar.Show()
			downloadRXSMUpdate(func(i int, text string) {
				sd.updateProgressBar.SetFormat(text + " - %p%")
				if i == -1 {
					success = false
					return
				}
				sd.updateProgressBar.SetValue(i)
			})
			if !success {
				return
			}
			sd.updateButton.SetEnabled(false)
			sd.updateButton.SetText("Updated (Restart Required)")
			sd.updateProgressBar.Hide()
		}()
	}
}

func (sd *SettingsDialog) onOkButtonClicked(bool) {
	sd.syncBackFields()
	sd.Accept()
}

func (sd *SettingsDialog) onCancelButtonClicked(bool) {
	sd.Reject()
}

// end SettingsDialog
