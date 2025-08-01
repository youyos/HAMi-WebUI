VERSION?=latest
DOCKER_IMAGE=projecthami/hami-webui-fe
OUT=./dist
PROJECT_NAME?=test-project

# 按项目最小化构建
ROUTE_FILE=packages/web/src/router/index.js
PROJECT_PATH=packages/web/projects/
DISABLED_PROJECTS?=""

.PHONY: install-modules
install-modules:
	pnpm install

.PHONY: build-all
build-all: install-modules build-bff build-web

.PHONY: build-bff
build-bff:
	pnpm run build

.PHONY: build-web
build-web:
	cd packages/web && pnpm run build

.PHONY: start-dev
start-dev: install-modules start-bff start-web


.PHONY: start-bff
start-bff:
	pnpm run start:dev &

.PHONY: start-web
start-web:
	cd packages/web && pnpm run start:dev

.PHONY: start-prod
start-prod:
	pnpm run start:prod

.PHONY: build-image
build-image:
	nerdctl -nk8s.io build --platform linux/amd64 -t ${DOCKER_IMAGE}:${VERSION} .

.PHONY: push-image
push-image:
	nerdctl -nk8s.io push ${DOCKER_IMAGE}:${VERSION}

.PHONY: release
release: build-image push-image
