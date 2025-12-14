@echo off
rem make protofile (Golang用)
rem *
rem あらかじめ全てのProto設定ファイルのimport文の下に
rem option go_package = ".;proto";
rem を追加すること
rem *
setlocal
protoc --go_out=. dwango\nicolive\chat\data\*.proto
protoc --go_out=. dwango\nicolive\chat\data\atoms\*.proto
protoc --go_out=. dwango\nicolive\chat\service\edge\*.proto
endlocal
