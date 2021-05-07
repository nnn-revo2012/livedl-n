#cd src
#go clean -cache
set-item env:GO111MODULE -value on
set-item env:GOARCH -value 386
set-item env:CGO_ENABLED -value 1
rm ..\livedl.x86.exe
go build -o ..\livedl.x86.exe .\livedl.go
#cd ..\
