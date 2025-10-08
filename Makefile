VERSION=v0.4.4

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}
