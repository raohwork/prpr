FROM golang AS builder
MAINTAINER ronmi.ren@gmail.com

ADD *.go go.mod go.sum /prpr/
WORKDIR /prpr
RUN go mod download && go build -v

FROM debian:unstable-slim
RUN apt-get update \
 && apt-get install --no-install-recommends -y firefox ca-certificates \
 && apt-get clean -y \
 && rm -fr /var/lib/apt/lists/*

ADD .mozilla /profile_tmpl
COPY --from=builder /prpr/prpr /usr/local/bin/

CMD /usr/local/bin/prpr
