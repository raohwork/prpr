FROM golang
MAINTAINER ronmi.ren@gmail.com

RUN apt-get update \
 && apt-get install -y firefox-esr \
 && apt-get clean -y \
 && rm -fr /var/lib/apt/lists/*

ADD *.go /go/src/github.com/raohwork/prpr/

WORKDIR /go/src/github.com/raohwork/prpr
RUN go get -v && go install github.com/raohwork/prpr

CMD /go/bin/prpr
