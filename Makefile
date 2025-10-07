VERSION=v0.4.0

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}
