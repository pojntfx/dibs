# syntax=docker/dockerfile:experimental
FROM --platform=$TARGETPLATFORM golang:1.13.5-buster AS build
WORKDIR /app
ARG TARGETPLATFORM

RUN go get github.com/magefile/mage
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./magefile.go ./magefile.go
COPY ./main.go ./main.go
COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN mage build

FROM --platform=$TARGETPLATFORM debian:buster-slim
ARG TARGETPLATFORM

COPY --from=build /app/.bin/godibs-* /usr/local/bin/godibs
EXPOSE 25000
CMD /usr/local/bin/godibs server
