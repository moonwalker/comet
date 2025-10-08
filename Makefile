VERSION=v0.4.6

# Bump version, commit, tag, and release
bump:
	@if git diff-index --quiet HEAD --; then \
		echo "❌ No changes to commit. Update VERSION in Makefile first."; \
		exit 1; \
	fi
	@echo "📦 Bumping to ${VERSION}..."
	git add -A
	git commit -m "chore: bump version to ${VERSION}"
	git tag -a ${VERSION} -m "Release ${VERSION}"
	git push && git push origin ${VERSION}
	@echo "✅ Released ${VERSION} - GitHub Action will build and publish"
