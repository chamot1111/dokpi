FROM debian:wheezy

# docker file for unit test
# docker build -t dokpi ./ && clear && docker run dokpi
# docker run -it --entrypoint=/bin/bash -i dokpi

# gcc for cgo
RUN apt-get update && apt-get install -y \
		gcc libc6-dev make\
		git-core\
		curl\
		wget\
		--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.5
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA1 5817fa4b2252afdb02e11e8b9dc1d9173ef3bd5a

ENV GITRECEIVE_URL https://raw.github.com/progrium/gitreceive/master/gitreceive

# git disable ssl verification
RUN git config --global http.sslVerify false

# WARNING: SSL verification canceled
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -k -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA1  golang.tar.gz" | sha1sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN curl -fsSL "$GITRECEIVE_URL" -k -o /usr/local/bin/gitreceive \
	&& chmod -R 777 /usr/local/bin/gitreceive \
	&& gitreceive init

# git test repo
ENV TEST_REPO /tests-repo
RUN mkdir -p "$TEST_REPO"
WORKDIR /tests-repo

RUN mkdir hello
WORKDIR /tests-repo/hello
ADD ./test_assets/test_buildpack/ ./
RUN touch unittest
RUN git init
RUN git add install detect deploy build unittest
RUN git commit -a -m "initial commit"

# add test buildpack
ENV PLUGINFOLDER /home/git/.dokpi/plugins
RUN mkdir -p "$PLUGINFOLDER/test" && chmod -R 777 "$PLUGINFOLDER"
ADD ./test/test_buildpack/ /home/git/.dokpi/plugins/test/

ADD ./ /dokpi/
WORKDIR /dokpi/
ENTRYPOINT ["go", "test", "./..."]
