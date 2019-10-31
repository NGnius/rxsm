// Created by NGnius 2019-10-31 sPoOKy

package main

import (
  "strings"
  "strconv"
)

// filter functions used in main-display.go
// these are all of the follwing form
// func myFunc(search string, s Save) bool {/* match code */}
// ie they take two parameters -- string and Save -- and return a boolean

func isAnyMatch(search string, s Save) bool {
  if isNameMatch(search, s) ||
    isCreatorMatch(search, s) ||
    isDescriptionMatch(search, s) ||
    isIDMatch(search, s) {
      return true
  }
  return false
}

func isNameMatch(search string, s Save) bool {
  if strings.Contains(strings.ToLower(s.Data.Name), strings.ToLower(search)) {
    return true
  }
  return false
}

func isCreatorMatch(search string, s Save) bool {
  if strings.Contains(strings.ToLower(s.Data.Creator), strings.ToLower(search)) {
    return true
  }
  return false
}

func isDescriptionMatch(search string, s Save) bool {
  if strings.Contains(strings.ToLower(s.Data.Description), strings.ToLower(search)) {
    return true
  }
  return false
}

func isIDMatch(search string, s Save) bool {
  stringID := strconv.Itoa(s.Data.Id)
  if strings.Contains(stringID, search) {
    return true
  }
  return false
}
