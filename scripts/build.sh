#!/bin/bash
rm -rf ./deploy/linux
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
cp ./resources/default_save/** ./deploy/linux/default_save
exit 0
