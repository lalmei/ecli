GO_BUILD64=GOARCH=amd64 go build
BIN=ecli

all: build

build:
	${GO_BUILD64} -o ${BIN}

darwin:
	@GOOS=darwin ${GO_BUILD64} -v -o ${BIN}

install:
	go install

clean:
	rm -f ecli ecli.xz

# Optional deploy rules
-include deploy.mk
