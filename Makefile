build:
	go build ./...

test:
	go test ./...

vet:
	go vet ./...

cover:
	go test ./... -coverprofile=c.out
	go tool cover -html="c.out"

lint:
	golangci-lint run

release:
ifndef v
	$(error Usage: make release v=0.2.0)
endif
	@echo "Releasing v$(v)..."
	git add -A
	git commit -m "chore: release v$(v)" --allow-empty
	git tag v$(v)
	git push
	git push origin v$(v)
	gh release create v$(v) --generate-notes --latest
	@echo "Done. v$(v) released."

version:
	@git describe --tags --abbrev=0 2>/dev/null || echo 'no tags'