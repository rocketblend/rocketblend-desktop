ARG BASE_IMAGE=ghcr.io/rocketblend/cross-wails:v2.8.0

FROM ${BASE_IMAGE} as builder

# RUN apt-get install -y osslsigncode

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG BUILD_TYPE=debug
ARG RELEASE_TAG=v0.0.0-0-g000000
ARG BUILD_TIMESTAMP=NOW
ARG COMMIT_SHA=docker
ARG BUILD_LINK=http://docker.local

RUN go run mage.go build ${BUILD_TYPE} ${RELEASE_TAG} ${BUILD_TIMESTAMP} ${COMMIT_SHA} ${BUILD_LINK}

RUN go test -v ./...

ENTRYPOINT [ "/bin/bash" ]

#############################################################

FROM ${BASE_IMAGE}

COPY --from=builder /usr/src/app/build/bin /out

ENTRYPOINT [ "sh", "-c" ]
CMD [ "cp -r /out/. /artifacts/" ]