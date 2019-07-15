.DEFAULT_GOAL := build

GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)

BINARY=tldr
LD_FLAGS += -s -w

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

release: clean fetch ## Generate a release, but don't publish to GitHub.
	$(GOPATH)/bin/goreleaser --skip-validate --skip-publish

publish: clean fetch ## Generate a release, and publish to GitHub.
	$(GOPATH)/bin/goreleaser

snapshot: clean fetch ## Generate a snapshot release.
	$(GOPATH)/bin/goreleaser --snapshot --skip-validate --skip-publish

fetch: ## Fetches the necessary dependencies to build.
	test -f $(GOPATH)/bin/goreleaser || go get -u -v github.com/goreleaser/goreleaser

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "dist/" "${BINARY}"

compress:
	(which upx > /dev/null && find dist/*/* | xargs -I{} -n1 -P 4 upx --best "{}") || echo "not using upx for binary compression"

build: fetch ## Builds the application.
	go build -ldflags "${LD_FLAGS}" -i -v -o ${BINARY}