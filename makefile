NAME := omxplayer-microservice
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}

COMMIT_HASH := $(shell git rev-parse --short HEAD)
TAG := $(shell git rev-parse --short HEAD)
ifneq ($(shell git describe --exact-match --tags HEAD 2> /dev/null),)
	TAG = $(shell git describe --exact-match --tags HEAD)
endif

PRD_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+"
DEV_TAG_REGEX := "v[0-9]+\.[0-9]+\.[0-9]+-.+"

# go stuff
PKG_LIST := $(shell go list ${PKG}/...)

.PHONY: all deps build test test-cov clean

all: clean build

test:
	@go test -coverprofile=coverage.txt -covermode=atomic ${PKG_LIST}

lint:
	@golangci-lint run --tests=false

deps:
	@echo Downloading dependencies...
	@go mod download

build: deps
	@mkdir -p dist
	@mkdir -p dist/files

	@echo
	@echo Building for linux-arm...
	@env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -v -o ./dist/${NAME} ${PKG}

	@echo Building deployment tarball
	@cp -r web ./dist/files/
	@cp -r static ./dist/files/
	@cd ./dist && tar -czf ${NAME}.tar.gz ${NAME} files

	@echo
	@echo Build output is located in ./dist/.

deploy: clean build
ifneq ($(shell echo ${TAG} | grep -x -E ${DEV_TAG_REGEX}),)
	@echo Getting current doc revision
	@$(eval rev=$(shell curl -s -n -X GET -u ${DB_USERNAME}:${DB_PASSWORD} "${DB_ADDRESS}/deployment-information/${NAME}" | cut -d, -f2 | cut -d\" -f4))

	@echo Pushing tar up to couch
	@curl -X PUT -u ${DB_USERNAME}:${DB_PASSWORD} -H "Content-Type: application/gzip" -H "If-Match: $(rev)" ${DB_ADDRESS}/deployment-information/${NAME}/development.tar.gz --data-binary @./dist/${NAME}.tar.gz
else ifneq ($(shell echo ${TAG} | grep -x -E ${PRD_TAG_REGEX}),)
	@echo Getting current doc revision
	@$(eval rev=$(shell curl -s -n -X GET -u ${DB_USERNAME}:${DB_PASSWORD} "${DB_ADDRESS}/deployment-information/${NAME}" | cut -d, -f2 | cut -d\" -f4))

	@echo Pushing tar up to couch
	@curl -X PUT -u ${DB_USERNAME}:${DB_PASSWORD} -H "Content-Type: application/gzip" -H "If-Match: $(rev)" ${DB_ADDRESS}/deployment-information/${NAME}/production.tar.gz --data-binary @./dist/${NAME}.tar.gz
endif

clean:
	@go clean
	@rm -rf ./dist/
