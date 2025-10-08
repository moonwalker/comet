VERSION=v0.4.6

tag:
	git tag -a ${VERSION} -m "Release ${VERSION}" && git push origin ${VERSION}
