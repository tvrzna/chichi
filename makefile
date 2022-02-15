DISTFILE=chichi
BUILD_VERSION=`git describe --tags`

test:
	@echo "Testing..."
	@go test -coverprofile cover.out ./...
	@echo "Done"

clean:
	@echo "Cleaning..."
	@rm -rf dist
	@echo "Done"

build:
	@echo "Building..."
	@mkdir -p dist
	@go build -o dist/${DISTFILE} -ldflags "-X github.com/tvrzna/chichi/src.buildVersion=${BUILD_VERSION}"
	@echo "Done"

install:
	@echo "Installing..."
	@install -DZs dist/${DISTFILE} -m 755 -t ${DESTDIR}/usr/bin
	@echo "Done"

uninstall:
	@echo "Uninstalling..."
	@rm -rf ${DESTDIR}/usr/bin/${DISTFILE}
	@echo "Done"
