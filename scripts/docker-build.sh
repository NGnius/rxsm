#!/bin/bash
rm -rf ./rxsm/deploy/linux
cd ./rxsm
$(go env GOPATH)/bin/qtdeploy -docker build linux
if [ $? -ne 0 ]
then
  exit $?
fi
# linux folder removal (patch)
if [ -d ./linux ]
then
  rmdir ./linux
fi
rm -f moc* # auto generated files
mkdir ./deploy/linux/default_save
cd ..
# build auto-updater binary
cd ./rxsm-updater/rxsm-updater
go build
cd ../..
cp ./rxsm-updater/rxsm-updater/rxsm-updater ./rxsm/deploy/linux/rxsm-updater
# copy resource files
cp ./resources/default_save/** ./rxsm/deploy/linux/default_save
mkdir ./rxsm/deploy/linux/icons
cp ./resources/exmods-icons/assets/*.svg ./rxsm/deploy/linux/icons
cp ./resources/exmods-icons/logos/rxsm-dual.svg ./rxsm/deploy/linux/icon.svg
cp ./resources/exmods-icons/logos/rxsm-dual-notext.svg ./rxsm/deploy/linux/icon-min.svg
cp ./resources/exmods-icons/assets/gear.svg ./rxsm/deploy/linux/settings.svg
cp README.md ./rxsm/deploy/linux/INFO.md
exit 0
