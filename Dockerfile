FROM golang
MAINTAINER ronmi.ren@gmail.com

RUN echo 'deb http://deb.debian.org/debian unstable main' > /etc/apt/sources.list.d/firefox.list \
 && apt-get update \
 && apt-get install -t unstable -y firefox \
 && apt-get clean -y \
 && rm -fr /etc/apt/sources.list.d/firefox.list /var/lib/apt/lists/*

ADD *.go /go/src/github.com/raohwork/prpr/
ADD .mozilla /root/.mozilla

WORKDIR /go/src/github.com/raohwork/prpr
RUN go get -v && go install github.com/raohwork/prpr

CMD /go/bin/prpr
