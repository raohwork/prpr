FROM golang
MAINTAINER ronmi.ren@gmail.com

RUN apt-get update \
 && apt-get install -y firefox-esr \
 && apt-get clean -y \
 && rm -fr /var/lib/apt/lists/*

ADD *.go /go/src/prpr/

WORKDIR /go/src/prpr
RUN go get -v && go install prpr

CMD /go/bin/prpr
