VERSION=v0.4.5

.PHONY: tag release-patch release-minor release-major

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}

# Bump patch version (0.0.x)
release-patch:
	@echo "Current version: ${VERSION}"
	@echo "Bumping patch version..."
	@NEW_VERSION=$$(echo ${VERSION} | awk -F. '{$$NF = $$NF + 1;} 1' OFS=.) && \
	sed -i.bak "s/VERSION=.*/VERSION=$$NEW_VERSION/" Makefile && \
	rm Makefile.bak && \
	git add Makefile && \
	git commit -m "chore: bump version to $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "version $$NEW_VERSION" && \
	git push && git push origin $$NEW_VERSION && \
	echo "✅ Released $$NEW_VERSION"

# Bump minor version (0.x.0)
release-minor:
	@echo "Current version: ${VERSION}"
	@echo "Bumping minor version..."
	@NEW_VERSION=$$(echo ${VERSION} | awk -F. '{$$2 = $$2 + 1; $$3 = 0;} 1' OFS=.) && \
	sed -i.bak "s/VERSION=.*/VERSION=$$NEW_VERSION/" Makefile && \
	rm Makefile.bak && \
	git add Makefile && \
	git commit -m "chore: bump version to $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "version $$NEW_VERSION" && \
	git push && git push origin $$NEW_VERSION && \
	echo "✅ Released $$NEW_VERSION"

# Bump major version (x.0.0)
release-major:
	@echo "Current version: ${VERSION}"
	@echo "Bumping major version..."
	@NEW_VERSION=$$(echo ${VERSION} | awk -F. '{$$1 = "v" (substr($$1,2) + 1); $$2 = 0; $$3 = 0;} 1' OFS=.) && \
	sed -i.bak "s/VERSION=.*/VERSION=$$NEW_VERSION/" Makefile && \
	rm Makefile.bak && \
	git add Makefile && \
	git commit -m "chore: bump version to $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "version $$NEW_VERSION" && \
	git push && git push origin $$NEW_VERSION && \
	echo "✅ Released $$NEW_VERSION"
