GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS )

BASE_PATH := $(shell pwd)
BUILD_PATH = $(BASE_PATH)/build
WEB_PATH=$(BASE_PATH)/frontend
ASSERT_PATH= $(BASE_PATH)/core/cmd/server/web/assets

CORE_PATH=$(BASE_PATH)/core
CORE_MAIN=$(CORE_PATH)/cmd/server/main.go
CORE_NAME=1panel-core

AGENT_PATH=$(BASE_PATH)/agent
AGENT_MAIN=$(AGENT_PATH)/cmd/server/main.go
AGENT_NAME=1panel-agent


clean_assets:
	rm -rf $(ASSERT_PATH)

upx_bin:
	upx $(BUILD_PATH)/$(CORE_NAME)
	upx $(BUILD_PATH)/$(AGENT_NAME)

build_frontend:
	cd $(WEB_PATH) && npm install && npm run build:pro

build_core_on_linux:
	cd $(CORE_PATH) \
	&& CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(CORE_NAME) $(CORE_MAIN)

build_agent_on_linux:
	cd $(AGENT_PATH) \
    && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(AGENT_NAME) $(AGENT_MAIN)

build_core_on_darwin:
	cd $(CORE_PATH) \
	&&  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -ldflags '-s -w'  -o $(BUILD_PATH)/$(CORE_NAME) $(CORE_MAIN)

build_agent_on_darwin:
	cd $(AGENT_PATH) \
    &&  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -ldflags '-s -w'  -o $(BUILD_PATH)/$(AGENT_NAME) $(AGENT_MAIN)

build_all: build_frontend build_core_on_linux build_agent_on_linux

build_on_local: clean_assets build_frontend build_core_on_darwin build_agent_on_darwin

# 发布相关
# VERSION   版本号，默认取 git tag，如 v2.0.0
# GOARCH    目标架构，默认取 go env GOARCH
# package_linux 生成 build/1panel-$(VERSION)-linux-$(GOARCH).tar.gz，
#               目录结构与 GitHub Releases 上的安装包一致。
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo v2.0.0)
PKG_NAME = 1panel-$(VERSION)-linux-$(GOARCH)
PKG_PATH = $(BUILD_PATH)/$(PKG_NAME)

# 面板内升级会把包内 1pctl 覆盖到 /usr/local/bin/1pctl，且只改写其中的 BASE_DIR 与 LANGUAGE，
# 因此包内 1pctl 的 ORIGINAL_VERSION 必须预先写入本次发布的版本号，
# 否则升级后 SystemVersion 会被回退成模板默认值。
package_linux: build_core_on_linux build_agent_on_linux
	rm -rf $(PKG_PATH)
	mkdir -p $(PKG_PATH)/initscript $(PKG_PATH)/lang
	cp $(BUILD_PATH)/$(CORE_NAME) $(PKG_PATH)/
	cp $(BUILD_PATH)/$(AGENT_NAME) $(PKG_PATH)/
	cp scripts/1pctl $(PKG_PATH)/
	sed -i "s#^ORIGINAL_VERSION=.*#ORIGINAL_VERSION=$(VERSION)#" $(PKG_PATH)/1pctl
	cp scripts/install.sh $(PKG_PATH)/
	cp scripts/lang/zh.sh $(PKG_PATH)/lang/
	cp scripts/lang/en.sh $(PKG_PATH)/lang/
	cp scripts/initscript/1panel-core.service $(PKG_PATH)/initscript/
	cp scripts/initscript/1panel-agent.service $(PKG_PATH)/initscript/
ifneq ($(wildcard assets/GeoIP.mmdb),)
	cp assets/GeoIP.mmdb $(PKG_PATH)/
endif
	chmod 755 $(PKG_PATH)/1pctl $(PKG_PATH)/install.sh
	cd $(BUILD_PATH) && tar czf $(PKG_NAME).tar.gz $(PKG_NAME)
	cd $(BUILD_PATH) && sha256sum $(PKG_NAME).tar.gz > checksums.txt
	@echo "package: $(BUILD_PATH)/$(PKG_NAME).tar.gz"

LAST_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "")

release_notes: package_linux
	@if [ -n "$(LAST_TAG)" ]; then \
		git log --no-merges --pretty=format:'- %s (%h)' $(LAST_TAG)..HEAD > $(BUILD_PATH)/$(PKG_NAME)-release-notes; \
	else \
		git log --no-merges --pretty=format:'- %s (%h)' -n 20 > $(BUILD_PATH)/$(PKG_NAME)-release-notes; \
	fi
	@echo "notes: $(BUILD_PATH)/$(PKG_NAME)-release-notes"
