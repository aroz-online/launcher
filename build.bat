echo "Building darwin"
set GOOS=darwin
set GOARCH=amd64

for %%I in (.) do SET EXENAME=%%~nxI

go build
MOVE "%EXENAME%" "%EXENAME%_darwin_amd64"

echo "Building linux"
set GOOS=linux
set GOARCH=amd64
go build
MOVE "%EXENAME%" "%EXENAME%_linux_amd64"

set GOOS=linux
set GOARCH=arm
go build
MOVE "%EXENAME%" "%EXENAME%_linux_arm"

set GOOS=linux
set GOARCH=arm64
go build
MOVE "%EXENAME%" "%EXENAME%_linux_arm64"

echo "Building windows"
set GOOS=windows
set GOARCH=amd64
go build

echo "Completed"
