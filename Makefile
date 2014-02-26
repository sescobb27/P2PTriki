GOPATH=$(shell pwd)
export GOPATH

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install
SRC=src/


TRIKI=triki
SERVER=server
CLIENT=client
$(shell cd $(SRC))

all:
	${GOBUILD} ${SERVER}.go
	${GOBUILD} ${CLIENT}.go
	${GOBUILD} ${TRIKI}
	${GOBUILD} ${TRIKI}

.PHONY: test open install

test:
	${GOTEST} ${TRIKI}/${SERVER}
	${GOTEST} ${TRIKI}/${CLIENT}