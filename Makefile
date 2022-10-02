.PHONY: build-agent
build-agent:
	go build -o ./build/gorong2-agent ./cmd/agent

.PHONY: build-server
build-server:
	go build -o ./build/gorong2-server ./cmd/server