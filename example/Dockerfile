# Use a multi-stage Docker build to keep our final image as small as possible.
# https://docs.docker.com/develop/develop-images/multistage-build/

FROM golang:alpine AS build

RUN apk add --no-cache git

# Instead of /go/src/app, which is typical for golang and the go-wrapper helper,
# we go straight to the proper path for this tool.  By doing this we avoid
# having to worry about creating the symbolic link, etc., and most go tools will
# *just work*.
WORKDIR /go/src/github.com/JaredReisinger/plugin-example
COPY . .

RUN go get -v .
RUN go build -v .

# Final stage: in a clean image, install the underlying command-line tool and
# then copy the newly-built plugin binary.

FROM alpine:latest

ENV HELM_VERSION="2.12.3"

LABEL maintainer="jaredreisinger@hotmail.com" \
      helm_version="${HELM_VERSION}"

RUN set -eux; \
    apk add --no-cache --virtual .build-deps \
        curl \
        ; \
    curl https://storage.googleapis.com/kubernetes-helm/helm-v${HELM_VERSION}-linux-amd64.tar.gz -o helm.tar.gz; \
    tar -zxvf helm.tar.gz; \
    for f in helm tiller; do install -t /usr/local/bin linux-amd64/$f; done; \
    rm -rf linux-amd64; \
    apk del .build-deps; \
    echo "DONE";

COPY --from=build /go/src/github.com/JaredReisinger/plugin-example/plugin-example /usr/local/bin/

CMD [ "/usr/local/bin/plugin-example" ]
