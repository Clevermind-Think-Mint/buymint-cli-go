# BuyMint CLI (GoLang Package)

This is a CLI client (or [GoLang](https://go.dev/) package) for [BuyMint](https://buy.bmint.it) SaaS.

Use this Go Package to integrate easily the Licesing features offered by [BuyMint](https://buy.bmint.it).

## As CLI

These are the commands you can use to interact with the BuyMint API with BuyMint CLI are available with the following command:

```sh
# To obtain a list of available commands
buymint-cli --help
# To obtain a list of available options for specific command
buymint-cli <command> --help
```

For example:

```sh
buymint-cli validate --help
```

## AS Package

Just use the package like this example:

```go
package "github.com/Clevermind-Think-Mint/buymint-cli-go"

# TODO

```

## Development

If you wish to collaborate with current project, you can initialize the project with the following steps:

```sh
# Esnure we are using Go modules
GO111MODULE=on # On Unix (otherwise, on Windows (Powershell), use "set GO111MODULE=on" or '$env:GO111MODULE="on"')
# Updating current module if needed
go mod tidy
# Launch the command to check if everything is all right
go run main.go --help
# Or, use any commands as, for example:
# Local license
go run main.go --debug --pretty validate -l ./test/assets/license.txt -p ./test/assets/public.key -m '{"agency": "A144109"}'
# Remote license (substite <serial> with desired one)
go run main.go --debug --pretty validate -l "https://buy.bmint.it/api/v1/service/microservice/licensor/license/<serial>" -m '{"agency": "A144109"}'
```

## Testing

If you wish to test your codebase launch the following command:

```sh
# Test the features
go test ./...
```

## Building

To build the command you can use the following commands.

### UNIX

On UNIX:

```sh
# Building Windows executables
GOOS="windows"; GOARCH="386"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/windows/x86/buymint-cli.exe
GOOS="windows"; GOARCH="amd64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/windows/amd64/buymint-cli.exe
# Building Linux executables
GOOS="linux"; GOARCH="386"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/x86/buymint-cli
GOOS="linux"; GOARCH="amd64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/amd64/buymint-cli
# Building Linux ARM executables
GOOS="linux"; GOARCH="arm"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/arm/buymint-cli
GOOS="linux"; GOARCH="arm64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/arm64/buymint-cli
# Building MAC executables
#GOOS="darwin"; GOARCH="386"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/darwin/x86/buymint-cli
GOOS="darwin"; GOARCH="amd64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/darwin/arm64/buymint-cli
```

### Windows (Powershell)

On Windows:

```sh
# Building Windows executables
$env:GOOS="windows"; $env:GOARCH="386"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/windows/x86/buymint-cli.exe
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/windows/amd64/buymint-cli.exe
# Building Linux executables
$env:GOOS="linux"; $env:GOARCH="386"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/x86/buymint-cli
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/amd64/buymint-cli
# Building Linux ARM executables
$env:GOOS="linux"; $env:GOARCH="arm"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/arm/buymint-cli
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/linux/arm64/buymint-cli
# Building MAC executables
#$env:GOOS="darwin";$env:GOARCH="386"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/darwin/x86/buymint-cli
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -ldflags "-X main.version=0.1.0 -X main.buildDate=01/05/2022" -o ./dist/darwin/arm64/buymint-cli
```