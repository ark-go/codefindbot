SHELL := /bin/bash
arklibgo := ~/ProjectsGo/arkAlias.sh
version = ~/ProjectsGo/arkAlias.sh getlastversion
.PHONY: check

.SILENT: build getlasttag buildzip buildwin


build:
	$(info +Компиляция Linux)
	@echo $$($(version))
	go build -ldflags "-s -w -X 'main.versionProg=$$($(version))'" -o ./bin/main/codefindbot cmd/main/main.go
buildzip:
	$(info +Компиляция с жатием)
	go build -ldflags "-s -w" -o ./bin/main/codefindbot cmd/main/main.go
buildwin:
	$(info +Компиляция windows)
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o ./bin/main/codefindbot.exe -tags static -ldflags "-s -w -X 'main.versionProg=$$($(version))'" cmd/main/main.go

run: build buildwin
	$(info +Запуск)
	./bin/main/codefindbot