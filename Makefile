VERSION := 0.1

build-metrics-server:
	@docker build -t velvetreactor/metrics:$(VERSION) .

.PHONY: all build-metrics-server