# vars
ORG=$(shell echo $(CIRCLE_PROJECT_USERNAME))
BRANCH=$(shell echo $(CIRCLE_BRANCH))
NAME=$(shell echo $(CIRCLE_PROJECT_REPONAME))

ifeq ($(NAME),)
NAME := $(shell basename "$(PWD)")
endif

ifeq ($(ORG),)
ORG=byuoitav
endif

ifeq ($(BRANCH),)
BRANCH:= $(shell git rev-parse --abbrev-ref HEAD)
endif

# go
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
VENDOR=gvt fetch -branch $(BRANCH)

all: deploy clean

build:
	env GOOS=linux GOARCH=arm $(GOBUILD) -o $(NAME) -v

test:
	$(GOTEST) -v -race $(go list ./... | grep -v /vendor/)

clean:
ifeq "$(BRANCH)" "master"
	$(eval BRANCH=development)
endif
	$(GOCLEAN)
	rm -f $(NAME)
	rm -f $(BRANCH).tar.gz
	rm -rf files/
	rm -rf vendor/
ifeq "$(BRANCH)" "development"
	$(eval BRANCH=master)
endif

deps:
ifneq "$(BRANCH)" "master"
	# put vendored packages in here
	gvt fetch -tag v3.3.10 github.com/labstack/echo
	$(VENDOR) github.com/byuoitav/common
endif
	$(GOGET) -d -v

deploy: $(NAME) version.txt
ifeq "$(BRANCH)" "master"
	$(eval BRANCH=development)
endif
	@echo Building deployment tarball
	@mkdir files
	@cp version.txt files/

	@tar -czf $(BRANCH).tar.gz $(NAME) files

	@echo Getting current doc revision
	$(eval rev=$(shell curl -s -n -X GET -u ${DB_USERNAME}:${DB_PASSWORD} "${DB_ADDRESS}/deployment-information/$(NAME)" | cut -d, -f2 | cut -d\" -f4))

	@echo Pushing zip up to couch
	@curl -X PUT -u ${DB_USERNAME}:${DB_PASSWORD} -H "Content-Type: application/gzip" -H "If-Match: $(rev)" ${DB_ADDRESS}/deployment-information/$(NAME)/$(BRANCH).tar.gz --data-binary @$(BRANCH).tar.gz
ifeq "$(BRANCH)" "development"
	$(eval BRANCH=master)
endif

### deps
$(NAME):
	$(MAKE) build
