build: build_amd64

build_amd64:
	@GOARCH=amd64 GOOS=linux go build -o bin/tbot-chatgpt-amd64 .

build_arm64:
	@GOARCH=arm64 GOOS=linux go build -o bin/tbot-chatgpt-arm64 .

build_darwin:
	@GOARCH=arm64 GOOS=darwin go build -o bin/tbot-chatgpt-arm64 .

all: build_amd64 build_arm64 build_darwin

clean:
	@rm -rf bin

docker_image:
	@docker build --tag tbot-chatgpt .
	docker save tbot-chatgpt > tbot-chatgpt.tar

.PHONY: build build_amd64 build_arm64 build_darwin clean docker_image
