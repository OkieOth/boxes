.PHONY: test build

VERSION = $(shell grep "const Version =" cmd/sub/version.go | grep "const Version =" | sed -e 's-.*= "--' -e 's-".*--')
UI_VERSION = $(shell cat ui/ui-version.txt)
CURRENT_DIR = $(shell pwd)
CURRENT_USER = $(shell id -u)
CURRENT_GROUP = $(shell id -g)

build:
	go build -o build/draw -ldflags "-s -w" cmd/main.go

build-docker:
	docker build -f Dockerfile.release -t ghcr.io/okieoth/boxes:$(VERSION) .

build-docker-ui:
	docker build -f ui/Dockerfile.ui -t ghcr.io/okieoth/boxes.ui:$(VERSION) .

build-wasm:
	rm -f ui/wasm/*.*
	GOOS=js GOARCH=wasm go build -o ui/wasm/boxes_$(VERSION).wasm wasm/main.go
	sed -i -e 's-window\.renderingVersion = \".*\";-window.renderingVersion = "v$(VERSION)";-' ui/html/index.html
	sed -i -e 's-\"/wasm/boxes_.*\.wasm\"-"/wasm/boxes_$(VERSION).wasm"-' ui/html/js/main.js

docker-push:
	docker push ghcr.io/okieoth/boxes:$(VERSION)

docker-ui-push:
	docker push ghcr.io/okieoth/boxes.ui:$(VERSION)

generate-all:
	bash -c scripts/generateAll.sh

run-ui-docker:
	docker run -p 8081:80 -d --rm ghcr.io/okieoth/boxes.ui:$(VERSION)

test:
	go test --cover ./... && echo ":)" || echo ":-/"
