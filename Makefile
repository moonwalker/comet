VERSION=v0.2.4

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}
