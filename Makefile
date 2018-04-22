PACKAGE_OS := linux darwin windows
PACKAGE_ARCH := amd64 386

build-image:
	docker build -t uphy/drone-util .

build:
	go build .

test:
	go test ./...

package: clean
	go get github.com/mitchellh/gox
	gox -os="$(PACKAGE_OS)" -arch="$(PACKAGE_ARCH)" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}"
	mkdir dist
	ls -1 build | xargs -I% cp config-example.yml README.md build/%/
	ls -1 build | xargs -I% tar zcf "dist/%.tar.gz" -C build "%"

clean:
	rm -rf build dist