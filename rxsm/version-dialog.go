// Created 2019-09-03 by NGnius

package main

import (
  "log"
  "strconv"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
  "gopkg.in/src-d/go-git.v4/plumbing/object"
  configlib "gopkg.in/src-d/go-git.v4/config"
  "github.com/therecipe/qt/widgets"
  "github.com/therecipe/qt/gui"
)

// NOTE: all "big" operations by go-git are slow so they are run in seperate goroutines

// start VersionDialog

type VersionDialog struct {
  widgets.QDialog
  saveVersioner ISaveVersioner
  isDetached bool
  infoLabel *widgets.QLabel
  settingsAutoLabel *widgets.QLabel
  settingsAutoField *widgets.QLineEdit
  treeView *widgets.QTreeWidget
  checkoutButton *widgets.QPushButton
  newVersionButton *widgets.QPushButton
  newBranchButton *widgets.QPushButton
  deleteBranchButton *widgets.QPushButton
  fillerLabel *widgets.QLabel
  closeButton *widgets.QPushButton
}

// NewVersionDialog(parent *widgets.QWidget, flags int) is automatically generated

func (vd *VersionDialog) OpenVersionDialog(saveVersioner ISaveVersioner) (int){
  vd.saveVersioner = saveVersioner
  vd.__init_display()
  vd.Open()
  return vd.Result()
}

func (vd *VersionDialog) __init_display() {
  vd.infoLabel = widgets.NewQLabel2("<b>Versions</b> <br/>Automatic snapshots are stopped while this menu is open", nil, 0)
  vd.infoLabel.SetTextFormat(1)
  vd.settingsAutoLabel = widgets.NewQLabel2("Take a snapshot every (seconds)", nil, 0)
  vd.settingsAutoLabel.SetWordWrap(true)
  vd.settingsAutoField = widgets.NewQLineEdit(nil)
  intValidator := gui.NewQIntValidator(nil)
  intValidator.SetBottom(0)
  vd.settingsAutoField.SetValidator(intValidator)
  vd.settingsAutoLabel.SetBuddy(vd.settingsAutoField)

  vd.treeView = widgets.NewQTreeWidget(nil)
  vd.treeView.SetHeaderLabels([]string{"Snapshot", "Hash"})

  vd.checkoutButton = widgets.NewQPushButton2("Go To", nil)
  vd.checkoutButton.ConnectClicked(vd.onCheckoutButtonClicked)
  vd.newVersionButton = widgets.NewQPushButton2("New Snapshot", nil)
  vd.newVersionButton.ConnectClicked(vd.onNewVersionButtonClicked)
  vd.newBranchButton = widgets.NewQPushButton2("New Branch", nil)
  vd.newBranchButton.ConnectClicked(vd.onNewBranchButtonClicked)
  vd.deleteBranchButton = widgets.NewQPushButton2("Delete Branch", nil)
  vd.deleteBranchButton.ConnectClicked(vd.onDeleteBranchButtonClicked)
  vd.deleteBranchButton.SetEnabled(false)

  vd.fillerLabel = widgets.NewQLabel2("", nil, 0)
  vd.fillerLabel.SetWordWrap(true)
  vd.closeButton = widgets.NewQPushButton2("Close", nil)
  vd.closeButton.ConnectClicked(vd.onCloseButtonClicked)

  headerLayout := widgets.NewQGridLayout2()
  headerLayout.AddWidget2(vd.infoLabel, 0, 0, 0)

  settingsLayout := widgets.NewQGridLayout2()
  settingsLayout.AddWidget2(vd.settingsAutoLabel, 0, 0, 0)
  settingsLayout.AddWidget2(vd.settingsAutoField, 0, 1, 0)

  versionLayout := widgets.NewQGridLayout2()
  versionLayout.AddWidget3(vd.treeView, 0, 0, 1, 2, 0)
  versionLayout.AddWidget2(vd.newVersionButton, 1, 0, 0)
  versionLayout.AddWidget2(vd.checkoutButton, 1, 1, 0)
  versionLayout.AddWidget2(vd.newBranchButton, 2, 0, 0)
  versionLayout.AddWidget2(vd.deleteBranchButton, 2, 1, 0)

  confirmLayout := widgets.NewQGridLayout2()
  confirmLayout.AddWidget3(vd.fillerLabel, 0, 0, 1, 4, 0)
  confirmLayout.AddWidget3(vd.closeButton, 0, 4, 1, 1, 0)

  masterLayout := widgets.NewQGridLayout2()
  masterLayout.AddLayout(headerLayout, 0, 0, 0)
  masterLayout.AddLayout(settingsLayout, 1, 0, 0)
  masterLayout.AddLayout(versionLayout, 2, 0, 0)
  masterLayout.AddLayout(confirmLayout, 3, 0, 0)

  vd.SetLayout(masterLayout)
  go func(){
    _, treeErr := makeTree(vd.saveVersioner.Repository(), vd.treeView)
    if treeErr != nil {
      log.Println("Error generating tree")
      log.Println(treeErr)
      return
    }
  }()
}

func (vd *VersionDialog) updateDetachedHeadWarning() {
  if vd.isDetached {
    vd.fillerLabel.SetText("New snapshots cannot be made when the current snapshot is not the latest in the branch!")
  } else {
    vd.fillerLabel.SetText("")
  }
}

func (vd *VersionDialog) onAutoFieldUpdate(value string) {
  // TODO: update config with new value
}

func (vd *VersionDialog) onCheckoutButtonClicked(bool) {
  selectedItem := vd.treeView.CurrentItem()
  if selectedItem == nil {
    log.Println("No tree item selected, ignoring checkout button click")
    return
  }
  checkoutOpts := &git.CheckoutOptions{Force: true}
  var debugHash string
  // TODO: alert user about possibility of losing work
  if vd.treeView.IndexOfTopLevelItem(selectedItem) == -1 {
    // commit/snapshot selected
    debugHash = "commit:"+selectedItem.Text(1)
    checkoutOpts.Hash = plumbing.NewHash(selectedItem.Text(1))
    if selectedItem.Parent().Child(selectedItem.Parent().ChildCount()-1).Text(1) != selectedItem.Text(1) {
      vd.isDetached = true // git head is detached; commits aren't possible in this state
    } else {
      vd.isDetached = false
    }
  } else {
    // branch selected
    debugHash = "branch:"+selectedItem.Text(1)
    checkoutOpts.Branch = plumbing.NewBranchReferenceName(selectedItem.Text(0))
    vd.treeView.SetCurrentItem(selectedItem.Child(selectedItem.ChildCount()-1))
    vd.isDetached = false
  }
  checkErr := vd.saveVersioner.Worktree().Checkout(checkoutOpts)
  if checkErr != nil {
    log.Println("Error during git checkout")
    log.Println(checkErr)
  }
  vd.updateDetachedHeadWarning()
  log.Println("Checked out version "+debugHash)
}

func (vd *VersionDialog) onNewVersionButtonClicked(bool) {
  if vd.isDetached {
    log.Println("New version cannot be created when HEAD detached, ignoring new version button click")
    return
  }
  selectedItem := vd.treeView.CurrentItem()
  if selectedItem == nil {
    log.Println("No tree item selected, ignoring new version button click")
    return
  }
  stagingDone := make(chan bool)
  go func(){
    vd.saveVersioner.StageAll()
    stagingDone <- true
  }()
  nameDialog := widgets.NewQInputDialog(nil, 0)
  var ok *bool // doesn't work properly
  name := nameDialog.GetText(nil, "Title", "Label", 0, "Text", ok, 0, 0)
  if name != "" {
    go func(){ // start goroutine func; this does not need to block UI actions
      <- stagingDone
      hashStr := vd.saveVersioner.PlainCommit(name) // slow :(
      if hashStr == "" {
        log.Println("Commit hash was empty, commit must have failed or been skipped")
        return
      }
      var selectedBranch *widgets.QTreeWidgetItem
      if vd.treeView.IndexOfTopLevelItem(selectedItem) == -1 {
        // commit/snapshot
        selectedBranch = selectedItem.Parent()
      } else {
        // branch
        selectedBranch = selectedItem
      }
      // add commit to tree widget
      commit, commitErr := vd.saveVersioner.Repository().CommitObject(plumbing.NewHash(hashStr))
      if commitErr != nil {
        log.Println("Error while retrieving commit object")
        log.Println(commitErr)
        return
      }
      selectedBranch.SetText(1, commit.Hash.String())
      newItem := widgets.NewQTreeWidgetItem7(nil, []string{commit.Message, commit.Hash.String()}, 0)
      selectedBranch.InsertChild(0, newItem)
      vd.treeView.SetCurrentItem(newItem)
      log.Println("New version of "+strconv.Itoa(vd.saveVersioner.Target().Data.Id)+" created")
    }() // end goroutine func
  } else {
    log.Println("New version cancelled")
  }
}

func (vd *VersionDialog) onNewBranchButtonClicked(bool) {
  selectedItem := vd.treeView.CurrentItem()
  if selectedItem == nil {
    log.Println("No tree item selected, ignoring new branch button click")
    return
  }
  nameDialog := widgets.NewQInputDialog(nil, 0)
  var ok *bool // doesn't work properly (sort of pointless, except compiler yells at me otherwise)
  name := nameDialog.GetText(nil, "Title", "Label", 0, "Text", ok, 0, 0)
  if name != "" {
    // create branch
    branchConf := &configlib.Branch{Name: name, Merge: plumbing.NewBranchReferenceName(name)}
    bErr := vd.saveVersioner.Repository().CreateBranch(branchConf)
    if bErr != nil {
      log.Println("Error creating branch")
      log.Println(bErr)
      return
    }
    // checkout branch
    checkoutOpts := &git.CheckoutOptions{Force: true, Create:true}
    checkoutOpts.Branch = plumbing.NewBranchReferenceName(name)
    checkoutOpts.Hash = plumbing.NewHash(selectedItem.Text(1))
    checkErr := vd.saveVersioner.Worktree().Checkout(checkoutOpts)
    if checkErr != nil {
      log.Println("Error checking out new branch")
      log.Println(checkErr)
      return
    }
    // add branch and commit to tree widget
    commit, commitErr := vd.saveVersioner.Repository().CommitObject(checkoutOpts.Hash)
    if commitErr != nil {
      log.Println("Error while retrieving commit object")
      log.Println(commitErr)
      return
    }
    newTopItem := widgets.NewQTreeWidgetItem4(nil, []string{name, checkoutOpts.Hash.String()}, 0)
    newItem := widgets.NewQTreeWidgetItem7(nil, []string{commit.Message, commit.Hash.String()}, 0)
    newTopItem.AddChild(newItem)
    vd.treeView.AddTopLevelItem(newTopItem)
    vd.treeView.SetCurrentItem(newItem)
    vd.isDetached = false
    vd.updateDetachedHeadWarning()
    log.Println("Created & checked out new branch "+name+" "+commit.Hash.String())
  } else {
    log.Println("New branch cancelled")
  }
}

func (vd *VersionDialog) onDeleteBranchButtonClicked(bool) {
  // TODO: fix branch not being properly deleted
  selectedItem := vd.treeView.CurrentItem()
  if selectedItem == nil {
    log.Println("No tree item selected, ignoring delete branch button click")
    return
  }
  if vd.treeView.IndexOfTopLevelItem(selectedItem) == -1 {
    log.Println("Commit tree item selected, ignoring delete branch button click")
    return
  }
  if selectedItem.Text(0) == "master" { // nothing good can come of this
    log.Println("Master branch selected, ignoring delete branch button click")
    return
  }
  delErr := vd.saveVersioner.Repository().DeleteBranch(selectedItem.Text(0))
  if delErr != nil {
    log.Println("Error deleting branch "+selectedItem.Text(0))
    log.Println(delErr)
    return
  }
  pruneErr := vd.saveVersioner.Repository().Prune(git.PruneOptions{
    Handler: func(hash plumbing.Hash) (error) {
      return vd.saveVersioner.Repository().DeleteObject(hash)
    }  })
  if pruneErr != nil {
    log.Println("Error pruning after deleting branch "+selectedItem.Text(0))
    log.Println(pruneErr)
    return
  }
  // remove branch (and children) from tree widget
  selectedItem.SetHidden(true) // deleting is too computationally hard
  log.Println("Deleted branch "+selectedItem.Text(0))
}

func (vd *VersionDialog) onCloseButtonClicked(bool) {
  vd.Accept()
}

// end VersionDialog

func makeTree(repo *git.Repository, treeWidget *widgets.QTreeWidget) (topItems []*widgets.QTreeWidgetItem, err error){
  // go through all branches' commits to count occurences of common parents
  head, err := repo.Head()
  var branchCommits []*object.Commit
  commitSet := NewCountSet()
  branchIter, branchErr := repo.Branches()
  if branchErr != nil {
    return topItems, branchErr
  }
  biterErr := branchIter.ForEach(func(branch *plumbing.Reference) (error) {
    commit, err := repo.CommitObject(branch.Hash())
    if err != nil {
      return err
    }
    branchCommits = append(branchCommits, commit)
    branchItem := widgets.NewQTreeWidgetItem4(nil, []string{branch.Name().Short(), branch.Hash().String()}, 0)
    treeWidget.AddTopLevelItem(branchItem)
    topItems = append(topItems, branchItem)
    for _, c := range getDumbAncestry(commit) {
      commitSet.Add(c.Hash.String())
    }
    return nil
  })
  if biterErr != nil {
    return topItems, biterErr
  }
  // go through all branches' commits again, stopping when a common parent is encountered
  for i, branchCom := range branchCommits {
    latestItem := widgets.NewQTreeWidgetItem7(nil, []string{branchCom.Message, branchCom.Hash.String()}, 0)
    topItems[i].AddChild(latestItem)
    if head.Hash().String() == branchCom.Hash.String() {
      treeWidget.SetCurrentItem(latestItem)
    }
    parentLoop: for _, parentCom := range getDumbAncestry(branchCom) {
      if commitSet.Count(parentCom.Hash.String()) > 1 && topItems[i].Data(0, 0).ToString() != "master" {
        break parentLoop
      }
      commitItem := widgets.NewQTreeWidgetItem7(nil, []string{parentCom.Message, parentCom.Hash.String()}, 0)
      topItems[i].AddChild(commitItem)
      if head.Hash().String() == parentCom.Hash.String() {
        treeWidget.SetCurrentItem(latestItem)
      }
    }
  }
  treeWidget.SetColumnCount(1)
  return
}

func getDumbAncestry(c *object.Commit) (ancestry []*object.Commit) {
  if len(c.ParentHashes) == 0 {
    return
  }
  parent, err := c.Parent(0) // it's dumb because it assumes only one parent
  // asexual reproduction ftw
  if err != nil {
    return
  }
  ancestry = append(ancestry, parent)
  ancestry = append(ancestry, getDumbAncestry(parent)...)
  return
}

// start CountSet (for checking if another object's hash already exists, and how many times)

type CountSet struct {
  items map[string]int
}

func NewCountSet() (*CountSet) {
  return &CountSet{items: map[string]int{}}
}

func (q *CountSet) Add(s string) {
  i, ok := q.items[s]
  if !ok {
    q.items[s] = 1
  } else {
    q.items[s] = i+1
  }
}

func (q *CountSet) Count(s string) (int) {
  i := q.items[s]
  return i
}

func (q *CountSet) Len() (int) {
  return len(q.items)
}

func (q *CountSet) Contains(s string) (bool) {
  _, ok := q.items[s]
  return ok
}

// end CountSet
