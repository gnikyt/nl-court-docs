BINARY_NAME=docket
PKG=github.com/gnikyt/nl-court-docs

all: clean build

build:
	go build -o ./dist/${BINARY_NAME} ${PKG}

clean:
	go clean
	rm ./dist/${BINARY_NAME} &2> /dev/null

docs:
	go doc ${PKG}

docs-all:
	go doc -all ${PKG}
