@echo off
echo Installing swag...
go install github.com/swaggo/swag/cmd/swag@latest
if %errorlevel% neq 0 (
    echo Failed to install swag
    exit /b %errorlevel%
)

echo Running swag init...
%USERPROFILE%\go\bin\swag init
if %errorlevel% neq 0 (
    echo Failed to run swag init. Trying just 'swag init'
    swag init
    if %errorlevel% neq 0 (
        echo Failed to run swag init
        exit /b %errorlevel%
    )
)

echo Running go mod tidy...
go mod tidy
