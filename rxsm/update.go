// Created 2019-09-12 by NGnius

package main

import (
  "net/http"
  "encoding/json"
  "io"
  "io/ioutil"
  "runtime"
  "bytes"
  "log"
)

const (
  DNT_ON = "1"
  DNT_OFF = "0"
)

var (
  ExtraHeader map[string][]string = make(map[string][]string)
  client http.Client = http.Client{}
)

type updateStruct struct {
  Status int `json:"status"`
  Reason string `json:"reason"`
  Url string `json:"url"`
  IsOutOfDate bool `json:"out-of-date"`
}

func CheckForUpdate(baseURL string, version string, platform string) (downloadURL string, isOutOfDate bool, ok bool) {
  body_map := make(map[string]string)
  body_map["version"] = version
  body_map["platform"] = platform
  body_bytes, marshalErr := json.Marshal(body_map)
  if marshalErr != nil {
    log.Println(marshalErr)
    return
  }
  req, _ := http.NewRequest("POST", baseURL+"/"+"update", bytes.NewReader(body_bytes))
  for key, elem := range ExtraHeader {
    req.Header[key] = elem
  }
  resp, httpErr := client.Do(req)
  if httpErr != nil || resp.StatusCode != 200 {
    log.Println(req.Header)
    log.Println(httpErr)
    if resp != nil {
      log.Println(resp.StatusCode)
    }
    return
  }
  defer resp.Body.Close()
  resp_struct := updateStruct{}
  resp_body_bytes, readAllErr := ioutil.ReadAll(resp.Body)
  if readAllErr != nil {
    log.Println(readAllErr)
    return
  }
  unmarshalErr := json.Unmarshal(resp_body_bytes, &resp_struct)
  if unmarshalErr != nil {
    log.Println(unmarshalErr)
    return
  }
  isOutOfDate = resp_struct.IsOutOfDate
  downloadURL = resp_struct.Url
  ok = true
  return
}

func SimpleCheckForUpdate(baseURL string, version string) (downloadURL string, isOutOfDate bool, ok bool) {
  return CheckForUpdate(baseURL, version, runtime.GOOS+"/"+runtime.GOARCH)
}

func DownloadUpdate(downloadURL string, dest io.Writer) (ok bool) {
  req, _ := http.NewRequest("GET", downloadURL, nil)
  for key, elem := range ExtraHeader {
    req.Header[key] = elem
  }
  resp, httpErr := client.Do(req)
  if httpErr != nil || resp.StatusCode != 200 {
    return
  }
  defer resp.Body.Close()
  _, copyErr := io.Copy(dest, resp.Body)
  if copyErr != nil {
    return
  }
  ok = true
  return
}
