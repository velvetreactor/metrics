release: git-tag-next build-metrics-server
	@docker image tag ghcr.io/velvetreactor/quant:$(VERSION) ghcr.io/velvetreactor/quant:latest
	@docker image push --all-tags ghcr.io/velvetreactor/quant

build-metrics-server:
	@docker buildx build --platform linux/amd64 -t ghcr.io/velvetreactor/quant:$(VERSION) .

git-tag-next:
	@git tag $(VERSION)
	@git push --tags

migrate-create:
	@migrate create -ext sql -dir migrations $(NAME)

migrate:
	@migrate -database $(DATABASE_URL) -path ./migrations up

.PHONY: all build-metrics-server release git-tag-next migrate migrate-create
