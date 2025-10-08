VERSION=v0.4.7

# Bump patch version, commit, tag, and release
bump:
	@echo "Current version: ${VERSION}"
	@NEW_VERSION=v0.4.7
	echo "Bumping to $$NEW_VERSION..." && \
	sed -i.bak "s/VERSION=v0.4.7
	rm -f Makefile.bak && \
	git add Makefile && \
	git commit -m "chore: bump version to $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "Release $$NEW_VERSION" && \
	git push && git push origin $$NEW_VERSION && \
	echo "✅ Released $$NEW_VERSION - GitHub Action will build and publish"

# Bump minor version (0.x.0)
bump-minor:
	@echo "Current version: ${VERSION}"
	@NEW_VERSION=v0.4.7
	echo "Bumping to $$NEW_VERSION..." && \
	sed -i.bak "s/VERSION=v0.4.7
	rm -f Makefile.bak && \
	git add Makefile && \
	git commit -m "chore: bump version to $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "Release $$NEW_VERSION" && \
	git push && git push origin $$NEW_VERSION && \
	echo "✅ Released $$NEW_VERSION - GitHub Action will build and publish"
