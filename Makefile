export GO111MODULE=on

COMMIT_HASH=$(shell git rev-parse --verify HEAD | cut -c 1-8)
BUILD_DATE=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_TAG=$(shell git describe --tags)

.PHONY: all
all: build test

.PHONY: build
build: clean
	@echo "compile"
	@go build -ldflags "-X 'main.AppVersion=`sh scripts/version.sh`'" cmd/adnet_server/main.go && \
	rm -rf build/adnet_server/* && \
	mkdir -p build/adnet_server/bin && mv main build/adnet_server/bin/adnet_server && \
	mkdir -p build/adnet_server/conf && cp conf/* build/adnet_server/conf && \
	cp scripts/adnet_op.sh build/adnet_server && \
	chmod +x scripts/check_adnet_server && \
	chmod +x scripts/reload && \
	cp scripts/check_adnet_server build/adnet_server && \
	cp scripts/reload build/adnet_server && \
	cp scripts/run build/adnet_server && \
	cp scripts/supervise.adnet_server build/adnet_server && \
	go build -ldflags "-X main.GitTag=$(GIT_TAG) -X main.BuildTime=$(BUILD_DATE) -X main.GitCommit=$(COMMIT_HASH)" ./cmd/hb_server/main.go && \
	rm -rf build/aladdin/* && \
	mkdir -p build/aladdin/bin && mv main build/aladdin/bin/aladdin_server && \
	mkdir -p build/aladdin/config && cp -r ./configs/* build/aladdin/config

.PHONY: mod
mod:
	go mod download

.PHONY: test
test:
	@echo "Run unit tests"
	- cd internal && go test -test.short -cover -gcflags=-l ./...

.PHONY: cover
cover:
	@echo "Run cover test"
	@go test -test.short -c -ldflags "-X main._VERSION_=$VERSION.${reversion}" -covermode=count \
	-coverpkg=$(shell echo `go list ./internal/...|grep -vE '^adnet$|/ad_server|/corsair_proto|/enum|/discovery|/netacuity|/watcher|/redis|/backender|/thrift_pool|/protobuf|/rbaserviceapi|/as_consul|/hb|/filter'`|sed 's/ /,/g') \
	./cmd/adnet_server -o adnet_server_cover && mkdir -p build/adnet_server/bin && \
	mv adnet_server_cover build/adnet_server/bin/adnet_server_cover && \
	go test -test.short -c -ldflags "-X main.GitTag=$(GIT_TAG) -X main.BuildTime=$(BUILD_DATE) -X main.GitCommit=$(COMMIT_HASH)" -covermode=count \
	-coverpkg=$(shell echo `go list ./internal/...|grep -vE '^aladdin$|/ad_server|/corsair_proto|/enum|/discovery|/netacuity|/watcher|/redis|/backender|/thrift_pool|/protobuf|/rbaserviceapi|/as_consul|/dnscache|/hb|/filter|/service|/watcher|/redis'`|sed 's/ /,/g') \
	./cmd/hb_server -o aladdin_server_cover && mkdir -p build/aladdin/bin && \
	mv aladdin_server_cover build/aladdin/bin/aladdin_server_cover

.PHONY: clean
clean:
	rm -rf build

.PHONY: gogoprotobuf
gogoprotobuf:
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/jsonpb
	go get github.com/gogo/protobuf/protoc-gen-gogo
	go get github.com/gogo/protobuf/gogoproto

.PHONY: docker
docker:
	@echo "Build aladdin docker image"
	docker build -t hub.mobvista.com/mae/aladdin:latest --target aladdin .
