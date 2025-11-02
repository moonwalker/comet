VERSION=v0.7.3

.PHONY: build
build:
	@echo "Building comet..."
	go build -o comet

.PHONY: release
release:
	@echo "Creating release ${VERSION}..."
	@echo "Make sure CHANGELOG.md is updated!"
	@read -p "Continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git add -A && \
		git commit -m "chore: release ${VERSION}" && \
		git tag -a ${VERSION} -m "Release ${VERSION}" && \
		git push && git push origin ${VERSION} && \
		echo "✅ Released ${VERSION} - GitHub Action will build and publish"; \
	else \
		echo "❌ Release cancelled"; \
	fi

.PHONY: website
website:
	@echo "Deploying website (triggers on push to main)..."
	@echo "Website changes will deploy automatically when pushed to main"
