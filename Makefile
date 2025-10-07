VERSION=v0.4.1

tag:
	@git tag -a ${VERSION} -m "version ${VERSION}" && git push origin ${VERSION}
