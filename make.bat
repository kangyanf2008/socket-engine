::编译linux版本
set GOOS=linux
set GOARCH=amd64
set GOHOSTOS=linux
go.exe build  -o bin/zore_client_A src/zore_client_A.go
go.exe build  -o bin/zore_client_B src/zore_client_B.go
go.exe build  -o bin/zore_server src/zore_server.go

::go modules proxy 设置 https://goproxy.io