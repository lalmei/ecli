.PHONY: build dist linux darwin windows buildall version cleardist clean

PROG=ecli
DISTDIR=dist
BINDIR=${DISTDIR}/${PROG}
#
VERSION=`git describe --tags --always`
BIN_VERSION=${PROG}-${VERSION}
BIN_DARWIN_AMD64=${BIN_VERSION}-darwin-amd64
BIN_FREEBSD_AMD64=${BIN_VERSION}-freebsd-amd64
BIN_LINUX_AMD64=${BIN_VERSION}-linux-amd64
BIN_WINDOWS_AMD64=${BIN_VERSION}-windows-amd64
#
BIN_BUILD64=GOARCH=amd64 go build
MAIN_CMD=github.com/keeneyetech/${PROG}

all: build

build: version
	@go build -v -o bin/${PROG}

install:
	go install

dist: distclean buildall zip

zip: linux darwin freebsd windows
	@rm -rf ${BINDIR}

linux:
	@cp bin/${BIN_VERSION}-linux* ${BINDIR}/${PROG} && \
		(cd ${DISTDIR} && zip -r ${BIN_LINUX_AMD64}.zip ${PROG})

darwin:
	@cp bin/${BIN_VERSION}-darwin* ${BINDIR}/${PROG} && \
		(cd ${DISTDIR} && zip -r ${BIN_DARWIN_AMD64}.zip ${PROG})

windows:
	@cp bin/${BIN_VERSION}-windows* ${BINDIR}/${PROG} && \
		(cd ${DISTDIR} && zip -r ${BIN_WINDOWS_AMD64}.zip ${PROG})

freebsd:
	@cp bin/${BIN_VERSION}-freebsd* ${BINDIR}/${PROG} && \
		(cd ${DISTDIR} && zip -r ${BIN_FREEBSD_AMD64}.zip ${PROG})

buildall: version
	@echo ">> DARWIN build"
	@GOOS=darwin ${BIN_BUILD64} -v -o bin/${BIN_DARWIN_AMD64} ${MAIN_CMD}
	@echo "\n>> FreeBSD build"
	@GOOS=freebsd ${BIN_BUILD64} -v -o bin/${BIN_FREEBSD_AMD64} ${MAIN_CMD}
	@echo "\n>> Linux build"
	@GOOS=linux ${BIN_BUILD64} -v -o bin/${BIN_LINUX_AMD64} ${MAIN_CMD}
	@echo "\n>> Windows build"
	@GOOS=windows ${BIN_BUILD64} -v -o bin/${BIN_WINDOWS_AMD64} ${MAIN_CMD}

version:
	@mkdir -p ${BINDIR}
	@./tools/version.sh

distclean: clean
	@rm -rf ${DISTDIR} && mkdir -p ${BINDIR}

clean:
	@rm -rf bin pkg *.xz ${DISTDIR}

# Optional deploy rules
-include deploy.mk
