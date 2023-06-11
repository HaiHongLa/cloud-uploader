@echo off

rem Build the project
go build -o file-storage.exe

:get-dependency
rem Download project dependencies
go mod download

echo Build completed.
