HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-backup
HELM_PLUGIN_NAME := "backup"
HAS_DEP := $(shell command -v dep;)
DEP_VERSION := v0.5.0
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: install
install: bootstrap build
	cp ${HELM_PLUGIN_NAME} $(HELM_PLUGIN_DIR)
	cp plugin.yaml $(HELM_PLUGIN_DIR)

.PHONY: hookInstall
hookInstall: bootstrap build

.PHONY: build
build:
	go build -o bin/${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./main.go

.PHONY: dist
dist:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-${HELM_PLUGIN_NAME}-linux-$(VERSION).tgz ${HELM_PLUGIN_NAME} README.md LICENSE.txt plugin.yaml
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-${HELM_PLUGIN_NAME}-macos-$(VERSION).tgz ${HELM_PLUGIN_NAME} README.md LICENSE.txt plugin.yaml
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${HELM_PLUGIN_NAME}.exe -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-${HELM_PLUGIN_NAME}-windows-$(VERSION).tgz ${HELM_PLUGIN_NAME}.exe README.md LICENSE.txt plugin.yaml
	rm ${HELM_PLUGIN_NAME}
	rm ${HELM_PLUGIN_NAME}.exe

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	wget -q -O $(GOPATH)/bin/dep https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-linux-amd64
	chmod +x $(GOPATH)/bin/dep
endif
	dep ensure
