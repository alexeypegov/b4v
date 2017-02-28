GO=$(shell which go)
GOINSTALL=$(GO) install
GOCLEAN=$(GO) clean
GOGET=$(GO) get
GOBUILD=$(GO) build

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=b4v
DOCKER_BIN=b4v_docker

DOCKER_LDFLAGS=--ldflags '-extldflags "-static"'
LDFLAGS=

.DEFAULT_GOAL: $(BINARY)
	
$(BINARY): $(SOURCES)
	$(GOBUILD) ${LDFLAGS} -o ${BINARY} main.go
	
$(DOCKER_BIN): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) ${DOCKER_LDFLAGS} -o ${DOCKER_BIN} main.go	
	
.PHONY: install get clean test all $(BINARY) $(DOCKER_BIN)
	
get:
	$(GOGET) github.com/BurntSushi/toml
	$(GOGET) github.com/bmizerany/pat
	$(GOGET) github.com/golang/glog
	$(GOGET) github.com/tylerb/graceful
	$(GOGET) github.com/urfave/negroni
	$(GOGET) github.com/boltdb/bolt
	
install:
	$(GOINSTALL) ${LDFLAGS} ./...
	
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${DOCKER_BIN} ] ; then rm ${DOCKER_BIN} ; fi
	
test:
	go test ./...
	
all: clean get $(BINARY) $(DOCKER_BIN)