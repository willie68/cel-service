@echo off
go build -ldflags="-s -w" -o service-cel-go.exe cmd/service/main.go