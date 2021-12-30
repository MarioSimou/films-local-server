FROM golang:1.16-buster

WORKDIR /opt
RUN apt-get update \
    && apt-get install zip wget -y \
    && wget -q https://github.com/aws/aws-sam-cli/releases/download/v1.36.0/aws-sam-cli-linux-x86_64.zip \
    && chmod +x aws-sam-cli-linux-x86_64.zip \
    && unzip aws-sam-cli-linux-x86_64.zip -d sam-installation \
    && ./sam-installation/install
    
WORKDIR /api
CMD [ "make", "test-api" ]