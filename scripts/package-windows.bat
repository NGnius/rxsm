ECHO OFF
cd .\rxsm && go mod download && go get -u -v github.com/therecipe/qt/cmd/qtdeploy && go get -u -v github.com/therecipe/qt/cmd/... && go mod vendor && git clone https://github.com/therecipe/env_windows_amd64_513.git vendor/github.com/therecipe/env_windows_amd64_513 && for /f %%v in ('go env GOPATH') do %%v\bin\qtdeploy build desktop && cd ..
xcopy .\resources\default_save .\rxsm\deploy\windows\default_save /I /Y