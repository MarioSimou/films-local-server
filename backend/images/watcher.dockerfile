FROM golang:1.16-buster

WORKDIR /opt
RUN apt-get update \
    && apt-get install make wget -y \
    && wget -q -O /opt/reflex_linux_amd64.tar.gz https://github.com/cespare/reflex/releases/download/v0.3.1/reflex_linux_amd64.tar.gz \
    && chmod +x /opt/reflex_linux_amd64.tar.gz \
    && tar -xf /opt/reflex_linux_amd64.tar.gz \
    && mv ./reflex_linux_amd64/reflex /go/bin/reflex \
    && rm -rf ./reflex_linux_amd64
    
WORKDIR /watcher
CMD [ "make", "watch-build-functions" ]