VERSION=v0.4.6

tag:
	@git diff-index --quiet HEAD || (echo "❌ You have uncommitted changes. Commit them first." && exit 1)
	git tag -a ${VERSION} -m "Release ${VERSION}" && git push origin ${VERSION}
