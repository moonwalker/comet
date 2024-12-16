VERSION=v0.3.1

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}
