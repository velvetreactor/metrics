VERSION := 0.2

release: build-metrics-server
	@docker image tag velvetreactor/metrics:$(VERSION) velvetreactor/metrics:latest

build-metrics-server:
	@docker build -t velvetreactor/metrics:$(VERSION) .

.PHONY: all build-metrics-server