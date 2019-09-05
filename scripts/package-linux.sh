#!/bin/bash
# dependencies
echo "Don't forget to update the version number in rxsm.go!"
cd ./rxsm
go mod download && go get -u -v github.com/therecipe/qt/cmd/qtdeploy && go get -u -v github.com/therecipe/qt/cmd/... && go mod vendor
git clone https://github.com/therecipe/env_linux_amd64_513.git vendor/github.com/therecipe/env_linux_amd64_513
cd ..
# build rxsm
./scripts/build.sh
