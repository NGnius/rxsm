ECHO OFF
cd .\rxsm && go mod download && go get -u -v github.com/therecipe/qt/cmd/qtdeploy && go get -u -v github.com/therecipe/qt/cmd/... && go mod vendor && git clone https://github.com/therecipe/env_windows_amd64_513.git vendor/github.com/therecipe/env_windows_amd64_513 && for /f %%v in ('go env GOPATH') do %%v\bin\qtdeploy build desktop && cd ..
cd .\rxsm-updater\rxsm-updater && go build && cd ..\.. && copy .\rxsm-updater\rxsm-updater\rxsm-updater.exe .\rxsm\deploy\windows\rxsm-updater.exe
xcopy .\resources\default_save .\rxsm\deploy\windows\default_save /I /Y
xcopy .\resources\exmods-icons\assets .\rxsm\deploy\windows\icons /I /Y
copy .\resources\exmods-icons\logos\rxsm-dual.svg .\rxsm\deploy\windows\icon.svg
copy .\resources\exmods-icons\logos\rxsm-dual-notext.svg .\rxsm\deploy\windows\icon-min.svg
copy .\resources\exmods-icons\assets\gear.svg .\rxsm\deploy\windows\settings.svg
copy README.md .\rxsm\deploy\windows\INFO.md
