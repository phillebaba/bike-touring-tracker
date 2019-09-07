.PHONY: app docker push

IMAGE ?= phillebaba/bike-touring-tracker
VERSION ?= $(shell git describe --tags --always --dirty)
TAG ?= $(VERSION)
PLATFORMS ?= "linux/amd64,linux/arm64,linux/arm"

app:
	go get github.com/golang/dep/cmd/dep
	dep ensure
	go get github.com/gobuffalo/packr/v2/packr2
	go get github.com/mitchellh/gox
	OSARCH=$$(echo $(PLATFORMS) | sed -e "s/,/ /g"); gox -osarch "$$OSARCH" -output="bike-touring-tracker-{{.Arch}}" ./cmd/server/

image: app
	docker run --rm --privileged multiarch/qemu-user-static:register
	docker buildx create --use --platform $(PLATFORMS)
	docker buildx build --platform $(PLATFORMS) -t $(IMAGE):$(TAG) --load .

push:	app
	docker run --rm --privileged multiarch/qemu-user-static:register
	docker buildx create --use --platform $(PLATFORMS)
	docker buildx build --platform $(PLATFORMS) -t $(IMAGE):$(TAG) --push .
