VERSION=v0.2.6

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}
