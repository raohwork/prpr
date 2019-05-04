FROM golang AS builder
MAINTAINER ronmi.ren@gmail.com

ADD *.go go.mod go.sum /prpr/
WORKDIR /prpr
RUN go mod download && go build -v



FROM golang AS downloader

RUN apt-get update \
 && apt-get install --no-install-recommends -y bzip2
RUN mkdir /firefox
WORKDIR /firefox
RUN wget -q -O - 'https://download.mozilla.org/?product=firefox-latest-ssl&os=linux64&lang=en-US' \
  | tar jxvf - --strip-components 1



FROM debian:stable-slim
RUN apt-get update \
 && apt-get upgrade -y \
 && apt-cache depends firefox-esr | grep Depends | cut -d ':' -f 2 \
  | xargs apt-get install -y --no-install-recommends ca-certificates \
 && apt-get clean -y \
 && rm -fr /var/lib/apt/lists/*

COPY --from=downloader /firefox /firefox
ADD .mozilla /profile_tmpl/
COPY --from=builder /prpr/prpr /usr/local/bin/
ADD start.sh /usr/local/bin/

CMD /usr/local/bin/start.sh
