# Useful Dockerfile with Go and Emboxen Framework.

FROM		debian:jessie
MAINTAINER	Tomoaki KAWADA <i@yvt.jp>

ENV DEBIAN_FRONTEND noninteractive

RUN	apt-get update \
&& apt-get -y dist-upgrade \
&& apt-get -y install \
golang ca-certificates mercurial git subversion nano apt-transport-https \
g++ \
&& rm -rf /var/lib/apt/lists/*

ENV GOPATH /go
WORKDIR $GOPATH

ADD . /go/src/github.com/yvt/emboxen-framework
RUN go get && go install github.com/yvt/emboxen-framework

