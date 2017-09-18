OUT_DIR = build
PACKAGE = k8s.io/event-adapter
PREFIX = gcr.io
GCP_PROJECT ?= google-containers
TAG = 1.0

PKG := $(shell find pkg/* -type f)
CMD := $(shell find cmd/* -type f)

build: build/adapter

build/adapter: adapter.go $(PKG) $(CMD)
	go build -a -o $(OUT_DIR)/adapter adapter.go

docker: build/adapter
	docker build --pull -t ${PREFIX}/${GCP_PROJECT}/event-adapter:$(TAG) .

push: docker
	gcloud docker -- push ${PREFIX}/${GCP_PROJECT}/event-adapter:$(TAG)

test:
	go test ./...

clean:
	rm -rf build apiserver.local.config
