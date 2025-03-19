BINARY_NAME=rapidusApp.exe
SHELL=cmd.exe

## build: builds all binaries
build:
# 	can use vendor below for local copy, but that is only really required in VS Code...
#    @go mod vendor
	@go build -o tmp/${BINARY_NAME} .
	@echo Rapidus built!

run: build
	@echo Staring Rapidus...
	@.\tmp\${BINARY_NAME} &
	@echo Rapidus started!

clean:
	@echo Cleaning...
	@go clean
	@del .\tmp\${BINARY_NAME}
	@echo Cleaned!

test:
	@echo Testing...
	@go test ./...
	@echo Done!

start: run

stop:
	@echo "Starting the front end..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo Stopped Rapidus

restart: stop start