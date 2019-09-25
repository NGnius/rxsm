#!/bin/bash
rm -rf ./rxsm/deploy/windows
cd ./rxsm
$(go env GOPATH)/bin/qtdeploy -docker build windows_64_shared
if [ $? -ne 0 ]
then
  exit $?
fi
# linux folder removal (patch)
if [ -d ./windows ]
then
  rmdir ./windows
fi
rm -f moc* # auto generated files
mkdir ./deploy/windows/default_save
cd ..
# build auto-updater binary
cd ./rxsm-updater/rxsm-updater
GOOS=windows GOARCH=amd64 go build
cd ../..
cp ./rxsm-updater/rxsm-updater/rxsm-updater.exe ./rxsm/deploy/windows/rxsm-updater.exe
# copy resource files
cp ./resources/default_save/** ./rxsm/deploy/windows/default_save
mkdir ./rxsm/deploy/windows/icons
cp ./resources/exmods-icons/assets/*.svg ./rxsm/deploy/windows/icons
cp ./resources/exmods-icons/logos/rxsm-dual.svg ./rxsm/deploy/windows/icon.svg
cp ./resources/exmods-icons/logos/rxsm-dual-notext.svg ./rxsm/deploy/windows/icon-min.svg
cp ./resources/exmods-icons/assets/gear.svg ./rxsm/deploy/windows/settings.svg
cp README.md ./rxsm/deploy/windows/INFO.md
exit 0
