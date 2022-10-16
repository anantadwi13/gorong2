.PHONY: build-agent
build-agent:
	go build -o ./build/gorong2-agent ./cmd/agent

.PHONY: build-server
build-server:
	go build -o ./build/gorong2-server ./cmd/server

.PHONY: docker-build-agent
docker-build-agent:
	docker build -t anantadwi13/gorong2-agent -f agent.dockerfile .

.PHONY: docker-build-server
docker-build-server:
	docker build -t anantadwi13/gorong2-server -f server.dockerfile .