VERSION := 0.4

release: build-metrics-server
	@docker image tag ghcr.io/velvetreactor/metrics:$(VERSION) ghcr.io/velvetreactor/metrics:latest
	@docker image push --all-tags ghcr.io/velvetreactor/metrics

build-metrics-server:
	@docker buildx build --platform linux/amd64 -t ghcr.io/velvetreactor/metrics:$(VERSION) .

.PHONY: all build-metrics-server
