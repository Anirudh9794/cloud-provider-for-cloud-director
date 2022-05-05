GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
GITROOT := $(shell git rev-parse --show-toplevel)
GO_CODE := $(shell ls go.mod go.sum **/*.go)
version := $(shell cat ${GITROOT}/release/version)

REGISTRY?="harbor-repo.vmware.com/vcloud"

.PHONY: build-within-docker

build-within-docker:
	mkdir -p /build/cloud-provider-for-cloud-director
	go build -ldflags "-X github.com/Anirudh9794/cloud-provider-for-cloud-director/version.Version=$(version)" -o /build/vcloud/cloud-provider-for-cloud-director cmd/ccm/main.go

ccm: $(GO_CODE)
	docker build -f Dockerfile . -t cloud-provider-for-cloud-director:$(version).latest
	docker tag cloud-provider-for-cloud-director:$(version).latest $(REGISTRY)/cloud-provider-for-cloud-director:$(version).latest
	# docker tag cloud-provider-for-cloud-director:$(version).latest $(REGISTRY)/cloud-provider-for-cloud-director:$(version).$$(docker images cloud-provider-for-cloud-director:$(version).latest -q)
	docker push $(REGISTRY)/cloud-provider-for-cloud-director:$(version).latest
	touch out/$@

all: ccm

test:
	go test -tags testing -v github.com/Anirudh9794/cloud-provider-for-cloud-director/pkg/vcdclient -cover -count=1
	go test -tags testing -v github.com/Anirudh9794/cloud-provider-for-cloud-director/pkg/config -cover -count=1

integration-test: test
	go test -tags="testing integration" -v github.com/Anirudh9794/cloud-provider-for-cloud-director/vcdclient -cover -count=1
