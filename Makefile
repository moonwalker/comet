VERSION=v0.2.2

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}
