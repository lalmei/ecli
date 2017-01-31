GO_BUILD64=GOARCH=amd64 go build
BIN=ecli
DARWIN_BIN=${BIN}-darwin

all: build

build:
	${GO_BUILD64} -o ${BIN}

darwin:
	@GOOS=darwin ${GO_BUILD64} -v -o ${DARWIN_BIN} && \
		xz ${DARWIN_BIN}

install:
	go install

clean:
	rm -f ecli* *.xz

# Optional deploy rules
-include deploy.mk
