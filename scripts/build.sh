#!/bin/bash
rm -rf ./rxsm/deploy/linux
cd ./rxsm
$(go env GOPATH)/bin/qtdeploy build desktop
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
cp ./resources/default_save/** ./rxsm/deploy/linux/default_save
cp ./resources/icon/rxsm-dual.svg ./rxsm/deploy/linux/icon.svg
cp ./resources/icon/rxsm-dual-notext.svg ./rxsm/deploy/linux/icon-min.svg
cp README.md ./rxsm/deploy/linux/INFO.md
exit 0
